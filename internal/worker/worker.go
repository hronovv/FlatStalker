package worker

import (
	"context"
	"log"
	"time"

	"flat-stalker/internal/repository"
	"flat-stalker/internal/source/kufar"
)

type Notifier interface {
	NotifyAd(ctx context.Context, chatID int64, text string) error
}

type Worker struct {
	interval time.Duration
	client   *kufar.Client
	listings *repository.Listings
	seen     *repository.SeenAds
	notifier Notifier
}

func New(
	interval time.Duration,
	listings *repository.Listings,
	seen *repository.SeenAds,
	notifier Notifier,
) *Worker {
	return &Worker{
		interval: interval,
		client:   kufar.NewClient(),
		listings: listings,
		seen:     seen,
		notifier: notifier,
	}
}

func (w *Worker) Start(ctx context.Context) {
	log.Printf("worker started, interval=%s", w.interval)
	w.runOnce(ctx)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("worker stopped")
			return
		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *Worker) runOnce(ctx context.Context) {
	watches, err := w.listings.ListWatches(ctx)
	if err != nil {
		log.Printf("worker: list watches: %v", err)
		return
	}
	if len(watches) == 0 {
		return
	}

	log.Printf("worker: checking %d watch(es)", len(watches))
	for _, watch := range watches {
		if err := w.checkWatch(ctx, watch); err != nil {
			log.Printf("worker: watch id=%d: %v", watch.ID, err)
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
	var freshIDs []int64
	for _, ad := range ads {
		if _, ok := existing[ad.AdID]; ok {
			continue
		}
		fresh = append(fresh, ad)
		freshIDs = append(freshIDs, ad.AdID)
	}
	if len(fresh) == 0 {
		return nil
	}

	for _, ad := range fresh {
		if err := w.notifier.NotifyAd(ctx, watch.ChatID, ad.FormatMessage()); err != nil {
			log.Printf("worker: notify chat_id=%d ad=%d: %v", watch.ChatID, ad.AdID, err)
			continue
		}
	}

	if err := w.seen.MarkMany(ctx, watch.ID, freshIDs); err != nil {
		return err
	}

	log.Printf("worker: watch id=%d notified %d new ad(s)", watch.ID, len(fresh))
	return nil
}
