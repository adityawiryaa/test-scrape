package usecases

import (
	"context"
	"log"
	"time"
)

func (q *queryUsecase) PollConfig(ctx context.Context) (int, error) {
	currentVersion := q.store.Version()

	cfg, changed, err := q.client.FetchConfig(ctx, currentVersion)
	if err != nil {
		return 0, err
	}

	if !changed {
		return 0, nil
	}

	q.store.Set(cfg)
	log.Printf("config updated to version %d", cfg.Version)

	return cfg.PollIntervalSeconds, nil
}

func (q *queryUsecase) StartPolling(ctx context.Context, initialInterval time.Duration, forwardFunc func(context.Context) error) {
	ticker := time.NewTicker(initialInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pollInterval, err := q.PollConfig(ctx)
			if err != nil {
				log.Printf("poll error: %v", err)
				continue
			}

			if pollInterval > 0 {
				ticker.Reset(time.Duration(pollInterval) * time.Second)

				if err := forwardFunc(ctx); err != nil {
					log.Printf("forward error: %v", err)
				}
			}
		}
	}
}
