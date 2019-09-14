package service

import (
	"context"
	"fmt"
	"net/http"
	"os"

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
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
	var addr string
	if fs.port == 0 {
		addr = fmt.Sprintf(":%s", os.Getenv("PORT"))
	} else {
		addr = fmt.Sprintf(":%d", fs.port)
	}
	return http.ListenAndServe(addr, nil)
}
