package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/mmcdole/gofeed"
	"sync"
)

var feedParser *gofeed.Parser

func init() {
	feedParser = gofeed.NewParser()
}
func MarshalFeed(url string) (*gofeed.Feed, error) {
	var feed *gofeed.Feed
	var err error
	if feed, err = feedParser.ParseURL(url); err != nil {
		return nil, err
	}
	return feed, err
}
func GetJSONString(obj interface{}, ch chan string, wg *sync.WaitGroup, acceptFields ...string) {
	toJson, err := json.Marshal(obj)
	if err != nil {
		ch <- ""
		wg.Done()
	}

	if len(acceptFields) == 0 {
		ch <- ""
		wg.Done()
	}

	toMap := map[string]interface{}{}
	json.Unmarshal([]byte(string(toJson)), &toMap)
	res := map[string]interface{}{}
	for _, field := range acceptFields {
		if v, found := toMap[field]; found {
			res[field] = v
		}
	}

	toJson, err = json.Marshal(res)
	if err != nil {
		ch <- ""
		wg.Done()
	}
	// push the population object down the channel
	ch <- string(toJson)
	// let the wait group know we finished
	wg.Done()
}
func RSSFeed(c *fiber.Ctx) error {
	q := c.Query("rssUrl")
	feed, err := MarshalFeed(q)
	if err != nil {
		panic(err)
	}
	var feeds []string
	// we create a buffered channel so writing to it won't block while we wait for the waitgroup to finish
	ch := make(chan string, len(feed.Items))
	wg := sync.WaitGroup{}
	for _, item := range feed.Items {
		wg.Add(1)
		go GetJSONString(item, ch, &wg, "content", "updatedParsed")
	}
	wg.Wait()
	close(ch)
	for jsonString := range ch {
		if len(jsonString) > 0 {
			feeds = append(feeds, jsonString)
		}
	}
	return c.JSON(fiber.Map{"items": feeds, "updated": feed.Updated})
}
