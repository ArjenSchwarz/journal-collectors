package main

import (
	"fmt"
	"time"

	"github.com/ArjenSchwarz/journal-collectors/config"
	collectors "github.com/ArjenSchwarz/journal-collectors/shared"
	"github.com/mmcdole/gofeed"
)

type configuration struct {
	CheckDuration    string
	Feeds            []string
	Dateformat       string
	IFTTTTriggerName string
	IFTTTKey         string
	KMS              bool
}

type parsedFeed struct {
	title string
	items []parsedItem
}

type parsedItem struct {
	title       string
	link        string
	description string
}

var settings configuration

func main() {
	err := config.ParseConfig(&settings)
	if err != nil {
		panic(err)
	}
	if settings.KMS {
		settings.IFTTTKey = collectors.DecryptKMS(settings.IFTTTKey)
	}
	feeds := make([]parsedFeed, 0, 1)
	for _, feed := range settings.Feeds {
		feeds = append(feeds, parseFeed(feed))
	}
	output := formatOutput(feeds)
	if output != "" {
		collectors.SendToIFTTT(settings.IFTTTKey, settings.IFTTTTriggerName, []string{output})
	}
}

func parseFeed(feedurl string) parsedFeed {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(feedurl)
	parsed := parsedFeed{title: feed.Title}
	checkDuration, err := time.ParseDuration(settings.CheckDuration)
	compareTime := time.Now().Add(-checkDuration)
	if err != nil {
		panic(err)
	}
	itemcollection := make([]parsedItem, 0, len(feed.Items))
	for _, item := range feed.Items {
		if compareTime.Before(*item.PublishedParsed) {
			itemparsed := parsedItem{
				title:       item.Title,
				link:        item.Link,
				description: item.Description}
			itemcollection = append(itemcollection, itemparsed)
		}
	}
	parsed.items = itemcollection
	return parsed
}

func formatOutput(feeds []parsedFeed) string {
	result := ""
	for _, feed := range feeds {
		if len(feed.items) > 0 {
			if result == "" {
				result = fmt.Sprintf("# Instapaper articles for %s", time.Now().Format(settings.Dateformat))
			}
			itemresult := fmt.Sprintf("<br><br>## %s", feed.title)
			for _, item := range feed.items {
				itemresult += fmt.Sprintf("<br> * [%s](%s)<br><br>%s", item.title, item.link, item.description)
			}
			result += itemresult
		}
	}
	return result
}
