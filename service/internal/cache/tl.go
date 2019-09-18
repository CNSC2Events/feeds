package cache

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/CNSC2Events/tlp"
	"github.com/rs/zerolog/log"
)

const (
	TLNET_SC2EVENTSURL = `https://liquipedia.net/starcraft2/api.php?action=parse&format=json&page=Liquipedia:Upcoming_and_ongoing_matches`
)

func fetchTL(ctx context.Context) (io.ReadCloser, error) {
	newContext, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(newContext, "GET", TLNET_SC2EVENTSURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func buildCache(ctx context.Context) error {
	r, err := fetchTL(ctx)
	if err != nil {
		return err
	}
	defer r.Close()
	p, err := tlp.NewTimelineParserFromReader(r)
	if err != nil {
		return err
	}
	if err := p.Parse(); err != nil {
		return err
	}
	if len(p.Events) == 0 {
		log.Error().Str("cache", "found no events").Send()
	}
	for _, e := range p.Events {
		data.Store(encodeFeedItemKey(p.RevID, e.VS.P1, e.VS.P2), e)
	}
	return nil
}

func encodeFeedItemKey(revID string, p1, p2 string) string {
	return fmt.Sprintf("%s-%s:%s", revID, p1, p2)
}

func GetAllMatches() []*tlp.Event {
	var matches []*tlp.Event
	data.Range(func(key, value interface{}) bool {
		log.Debug().Interface(key.(string), value)
		if m, ok := value.(*tlp.Event); ok {
			matches = append(matches, m)
		}
		return true
	})
	return matches
}
