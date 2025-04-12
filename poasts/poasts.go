package poasts

// The poasts package handles reading the posts from the RSS feed.
// It takes in

import (
	"iter"
	"time"

	"github.com/araddon/dateparse"
	"github.com/mmcdole/gofeed"
)

// A raw stream of posts
type Poasts struct {
	// The underlying URL of the RSS feed
	URL  string
	feed *gofeed.Feed
}

// A pointer to a blog poast.
type Poast struct {
	// The Title description of the Poast
	Title string

	// The short-form Description of the Poast
	Description string

	// The URL of the Poast
	URL string

	// The time of publication
	Published time.Time
}

// Create a new stream of posts.
// Accepts an RSS feed's URL or local path, and a SQLite path.
// Returns a *Poasts iterator
func New(rss string) (*Poasts, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rss)
	if err != nil {
		return nil, err
	}
	return &Poasts{URL: rss, feed: feed}, nil
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
				URL:         itm.GUID,
				Published:   published,
			}
			if !yield(&p, nil) {
				break
			}
		}
	}
}

// All of the Poasts from after the given publication time.
func (ps *Poasts) After(published *time.Time) iter.Seq2[*Poast, error] {
	if published == nil {
		return ps.All()
	}
	return func(yield func(*Poast, error) bool) {
		for p, err := range ps.All() {
			if err != nil {
				yield(nil, err)
				break
			}
			if p.Published.After(*published) && !yield(p, nil) {
				break
			}
		}
	}
}
