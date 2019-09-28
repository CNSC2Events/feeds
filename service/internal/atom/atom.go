// package atom parses the tl.net events to atom feeds
package atom

import (
	"bytes"
	"fmt"
	"strings"
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
		var dURL string
		if item.DetailURL != nil {
			dURL = item.DetailURL.String()
		}
		items = append(items, &feeds.Item{
			Title:       buildTitle(item),
			Link:        &feeds.Link{Href: dURL},
			Description: "",
			Author:      &feeds.Author{Name: "Nace Sc", Email: "scbizu@gmail.com"},
			Created:     now,
			Content:     buildHTMLContent(item),
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
	vs := strings.ReplaceAll(e.GetVersus(), "\n", "")
	if e.IsOnGoing {
		return fmt.Sprintf("[OnGoing][%s]%s", e.Series, vs)
	}
	return fmt.Sprintf("[%s]%s", e.Series, vs)
}

func buildHTMLContent(e *tlp.Event) string {
	tmpl := `<html>
	<body>
		<p>Series: <span id="serires">%s</span></p>
		<p>VS: <span id="vs">%s</span></p>
		<p>Fight At: <span id="start_at">%s</span></p>
	</body>
	</html>`

	vs := strings.ReplaceAll(e.GetVersus(), "\n", "")
	return fmt.Sprintf(tmpl, e.Series, vs, e.StartAt.String())
}
