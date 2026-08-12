package worker

import (
	"context"
	"log"
	"time"

	"flat-stalker/internal/repository"
	"flat-stalker/internal/source/kufar"
)

// SeedFromURL fetches the search snapshot and marks those ads as seen without notifying.
func SeedFromURL(ctx context.Context, client *kufar.Client, seen *repository.SeenAds, listingID int64, searchURL string) (int, error) {
	ads, err := client.FetchAds(searchURL)
	if err != nil {
		return 0, err
	}
	if len(ads) == 0 {
		return 0, nil
	}
	adIDs := make([]int64, len(ads))
	for i, ad := range ads {
		adIDs[i] = ad.AdID
	}
	if err := seen.MarkMany(ctx, listingID, adIDs); err != nil {
		return 0, err
	}
	return len(adIDs), nil
}

func SeedFromURLAsync(client *kufar.Client, seen *repository.SeenAds, listingID int64, searchURL string) {
	if client == nil || seen == nil {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		n, err := SeedFromURL(ctx, client, seen, listingID, searchURL)
		if err != nil {
			log.Printf("seed listing id=%d: %v", listingID, err)
			return
		}
		log.Printf("seed listing id=%d: %d ads", listingID, n)
	}()
}
