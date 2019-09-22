package cache

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var rdyChan chan struct{}

var data *sync.Map

// Register inits the feeds cache
func Init(ctx context.Context) {
	data = new(sync.Map)

	go refreshCache(ctx)
}

func GetReadySignal() {
	if getSize() > 0 {
		return
	}
	for {
		select {
		case <-rdyChan:
			return
		}
	}
}

func refreshCache(ctx context.Context) {

	if getSize() == 0 {
		refreshOnceIfEmpty(ctx)
		rdyChan <- struct{}{}
		return
	}

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
	if err := buildCache(ctx); err != nil {
		log.Warn().Err(err).Send()
	}
}

func getSize() int {

	var len int

	data.Range(func(k, v interface{}) bool {
		len++
		return true
	})

	return len
}
