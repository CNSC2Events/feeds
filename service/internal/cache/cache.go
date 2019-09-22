package cache

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var data *sync.Map

// Register inits the feeds cache
func Init(ctx context.Context) {
	data = new(sync.Map)

	// serverless app will not spawn up a goroutine
	if os.Getenv("SERVERLESS") == "" {
		go refreshCache(ctx)
	}
}

func refreshCache(ctx context.Context) {

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
