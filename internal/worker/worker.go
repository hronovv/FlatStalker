package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"flat-stalker/internal/plan"
	"flat-stalker/internal/repository"
	"flat-stalker/internal/source/kufar"
)

type Notifier interface {
	NotifyAd(ctx context.Context, chatID int64, text string) error
}

type Worker struct {
	intervals plan.Intervals
	client    *kufar.Client
	listings  *repository.Listings
	seen      *repository.SeenAds
	notifier  Notifier
}

func New(
	intervals plan.Intervals,
	listings *repository.Listings,
	seen *repository.SeenAds,
	notifier Notifier,
) *Worker {
	return &Worker{
		intervals: intervals,
		client:    kufar.NewClient(),
		listings:  listings,
		seen:      seen,
		notifier:  notifier,
	}
}

func (w *Worker) Start(ctx context.Context) {
	var wg sync.WaitGroup
	tiers := []struct {
		name     string
		interval time.Duration
	}{
		{plan.Pro, w.intervals.Pro},
		{plan.Plus, w.intervals.Plus},
		{plan.Free, w.intervals.Free},
	}

	for _, tier := range tiers {
		wg.Add(1)
		go func(name string, interval time.Duration) {
			defer wg.Done()
			w.runTier(ctx, name, interval)
		}(tier.name, tier.interval)
	}

	wg.Wait()
	log.Println("worker stopped")
}

func (w *Worker) runTier(ctx context.Context, userPlan string, interval time.Duration) {
	log.Printf("worker[%s] started, interval=%s", userPlan, interval)
	w.runOnce(ctx, userPlan)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker[%s] stopped", userPlan)
			return
		case <-ticker.C:
			w.runOnce(ctx, userPlan)
		}
	}
}

func (w *Worker) runOnce(ctx context.Context, userPlan string) {
	watches, err := w.listings.ListWatches(ctx, userPlan)
	if err != nil {
		log.Printf("worker[%s]: list watches: %v", userPlan, err)
		return
	}
	if len(watches) == 0 {
		return
	}

	log.Printf("worker[%s]: checking %d watch(es)", userPlan, len(watches))
	for _, watch := range watches {
		if err := w.checkWatch(ctx, watch); err != nil {
			log.Printf("worker[%s]: watch id=%d: %v", userPlan, watch.ID, err)
		}
	}
}

func (w *Worker) checkWatch(ctx context.Context, watch repository.Watch) error {
	ads, err := w.client.FetchAds(watch.URL)
	if err != nil {
		return err
	}
	if len(ads) == 0 {
		return nil
	}

	adIDs := make([]int64, len(ads))
	for i, ad := range ads {
		adIDs[i] = ad.AdID
	}

	count, err := w.seen.CountByListing(ctx, watch.ID)
	if err != nil {
		return err
	}

	// First poll: seed current ads without spamming the user.
	if count == 0 {
		if err := w.seen.MarkMany(ctx, watch.ID, adIDs); err != nil {
			return err
		}
		log.Printf("worker: watch id=%d first run, seeded %d ads", watch.ID, len(adIDs))
		return nil
	}

	existing, err := w.seen.ExistingIDs(ctx, watch.ID, adIDs)
	if err != nil {
		return err
	}

	var fresh []kufar.Ad
	for _, ad := range ads {
		if _, ok := existing[ad.AdID]; ok {
			continue
		}
		fresh = append(fresh, ad)
	}
	if len(fresh) == 0 {
		return nil
	}

	var notified []int64
	for _, ad := range fresh {
		if err := w.notifier.NotifyAd(ctx, watch.ChatID, ad.FormatMessage()); err != nil {
			log.Printf("worker: notify chat_id=%d ad=%d: %v", watch.ChatID, ad.AdID, err)
			continue
		}
		notified = append(notified, ad.AdID)
	}
	if len(notified) == 0 {
		return nil
	}

	if err := w.seen.MarkMany(ctx, watch.ID, notified); err != nil {
		return err
	}

	log.Printf("worker: watch id=%d notified %d new ad(s)", watch.ID, len(notified))
	return nil
}
