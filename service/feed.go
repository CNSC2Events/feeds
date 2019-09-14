package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/CNSC2Events/feeds/service/internal/atom"
	"github.com/CNSC2Events/feeds/service/internal/cache"
	"github.com/rs/zerolog/log"
)

type FeedService struct {
	port int32
}

func RegisterCache(ctx context.Context) {
	cache.Init(ctx)
}

func NewFeedService(port int32) FeedService {
	return FeedService{port: port}
}

func (fs FeedService) Serve(ctx context.Context) error {

	http.HandleFunc("/sc2", func(w http.ResponseWriter, r *http.Request) {
		data := cache.GetAllMatches()
		f := atom.Feeds{Items: data}
		ww, err := f.Writer()
		if err != nil {
			log.Error().Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
		}
		if _, err := w.Write(ww.Bytes()); err != nil {
			log.Error().Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", fs.port), nil)
}
