package main

import (
	"encoding/xml"
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

func fetchFeed(feedUrl string) (*RSSFeed, error) {
    req, err := http.NewRequest("GET", feedUrl, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("User-Agent", "Gator")

    client := &http.Client{}

    res, err := client.Do(req)
    if err != nil {
        return nil, err
    }

    defer res.Body.Close()

    data, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    rssFeed := &RSSFeed{}

    if err = xml.Unmarshal(data, rssFeed); err != nil {
        return nil, err
    }

    rssFeed.cleanUpTitleDesc()

    return rssFeed, nil
}

/*
This function use html.UnescapeString to clean up
the title and the description of the rss feed
*/
func (rssFeed *RSSFeed) cleanUpTitleDesc() {
    rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
    rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

    for i, item := range rssFeed.Channel.Item {
        item.Title = html.UnescapeString(item.Title)
        item.Description = html.UnescapeString(item.Description)
        rssFeed.Channel.Item[i] = item
    }
}

