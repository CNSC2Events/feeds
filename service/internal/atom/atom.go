// package atom parses the tl.net events to atom feeds
package atom

import (
	"bytes"
	"fmt"
	"time"

	"github.com/CNSC2Events/tlp"
	"github.com/gorilla/feeds"
)

type Feeds struct {
	Items []*tlp.Event
}

func (f Feeds) Writer() (*bytes.Buffer, error) {
	b := bytes.NewBufferString("")
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "tl.net feeds",
		Link:        &feeds.Link{Href: "https://rss.scnace.me/sc2"},
		Description: "tl.net living match feeds",
		Author:      &feeds.Author{Name: "Nace Sc", Email: "scbizu@gmail.com"},
		Created:     now,
	}

	var items []*feeds.Item

	for _, item := range f.Items {
		items = append(items, &feeds.Item{
			Title:       buildTitle(item),
			Link:        &feeds.Link{Href: item.DetailURL.String()},
			Description: "",
			Author:      &feeds.Author{Name: "Nace Sc", Email: "scbizu@gmail.com"},
			Created:     now,
		})
	}

	feed.Items = items

	atom, err := feed.ToAtom()
	if err != nil {
		return nil, err
	}
	b.WriteString(atom)
	return b, nil
}

func buildTitle(e *tlp.Event) string {
	return fmt.Sprintf("[%s]%s", e.Series, e.GetVersus())
}
