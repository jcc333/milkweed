package milkweed

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func main() {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("http://welcome-to-dudely-house.net/feed.xml")
	fmt.Println(feed.Title)
}
