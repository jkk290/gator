package feedapi

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedurl string) (*RSSFeed, error) {
	newRequest, err := http.NewRequestWithContext(ctx, "GET", feedurl, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error creating new request: %w", err)
	}
	newRequest.Header.Set("User-Agent", "gator")
	res, err := http.DefaultClient.Do(newRequest)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error getting response: %w", err)
	}
	defer res.Body.Close()

	resBodyByte, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error reading from response body: %w", err)
	}

	newFeed := RSSFeed{}
	if err := xml.Unmarshal(resBodyByte, &newFeed); err != nil {
		return &RSSFeed{}, fmt.Errorf("error converting from xml: %w", err)
	}
	newFeed.Channel.Title = html.UnescapeString(newFeed.Channel.Title)
	newFeed.Channel.Description = html.UnescapeString(newFeed.Channel.Description)
	for i, item := range newFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		newFeed.Channel.Item[i] = item
	}
	return &newFeed, nil
}
