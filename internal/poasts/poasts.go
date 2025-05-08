package poasts

// The poasts package handles reading the posts from the RSS feed.
// It takes in

import (
	"fmt"
	"iter"
	"log"
	"time"

	"github.com/araddon/dateparse"
	"github.com/mmcdole/gofeed"
)

// A raw stream of posts
type Poasts struct {
	// The underlying URL of the RSS feed
	URL    string
	feed   *gofeed.Feed
	logger *log.Logger
}

// A pointer to a blog poast.
type Poast struct {
	// The description of the Poast
	Title string

	// The short-form description of the Poast
	Description string

	// The GUID of the Poast
	GUID string

	// The time of publication
	Published time.Time
}

func (p *Poast) String() (skeet string, err error) {
	skeet = fmt.Sprintf("%s: %s\n%s", p.Title, p.GUID, p.Description)
	if len(skeet) > 300 {
		skeet = fmt.Sprintf("%s: %s", p.Title, p.GUID)
		if len(skeet) > 300 {
			skeet = fmt.Sprintf("%s", p.GUID)
			if len(skeet) > 300 {
				err = fmt.Errorf("Milkweed doesn't have a minifier yet, and this GUID is too long. Sorry")
			}
		}
	}
	return
}

// Create a new stream of posts.
// Accepts an RSS feed's URL or local path, and a SQLite path.
// Returns a *Poasts iterator
func New(rss string, logger *log.Logger) (*Poasts, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rss)
	if err != nil {
		return nil, err
	}
	logger.Println("Parsed RSS feed for '", feed.Title, "' at URL '", rss, "'")
	return &Poasts{URL: rss, feed: feed, logger: logger}, nil
}

// All of the `Poast`s for the, ah, `Poasts`.
func (ps *Poasts) All() iter.Seq2[*Poast, error] {
	return func(yield func(*Poast, error) bool) {
		// iterate each row of the RSS feed
		for _, itm := range ps.feed.Items {
			published, err := dateparse.ParseAny(itm.Published)
			if err != nil {
				yield(nil, err)
				break
			}
			p := Poast{
				Title:       itm.Title,
				Description: itm.Description,
				GUID:        itm.GUID,
				Published:   published,
			}
			if !yield(&p, nil) {
				break
			}
		}
	}
}
