package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")
	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()
	response, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	var rssFeed RSSFeed
	if err := xml.Unmarshal(response, &rssFeed); err != nil {
		return &RSSFeed{}, err
	}
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for index, content := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[index].Title = html.UnescapeString(content.Title)
		rssFeed.Channel.Item[index].Description = html.UnescapeString(content.Description)
	}
	return &rssFeed, nil

}
