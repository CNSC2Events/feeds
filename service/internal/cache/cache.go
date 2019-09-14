package cache

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var data *sync.Map

// Register inits the feeds cache
func Init(ctx context.Context) {
	data = new(sync.Map)

	go refreshCache(ctx)
}

func refreshCache(ctx context.Context) {

	refreshOnceIfEmpty(ctx)

	ticker := time.NewTicker(5 * time.Minute)

	for {
		select {
		case <-ticker.C:
			log.Debug().Msg("cache: start to rebuild")
			if err := buildCache(ctx); err != nil {
				log.Warn().Err(err).Send()
			}
		}
	}

}

func refreshOnceIfEmpty(ctx context.Context) {
	var len int

	data.Range(func(k, v interface{}) bool {
		len++
		return true
	})

	if len == 0 {
		if err := buildCache(ctx); err != nil {
			log.Warn().Err(err).Send()
		}
	}
}
