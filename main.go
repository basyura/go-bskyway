package main

import (
	"context"
	"fmt"

	"bskyway/config"
	"bskyway/session"

	"github.com/bluesky-social/indigo/api/bsky"
)

func main() {
	if err := doMain(); err != nil {
		fmt.Println(err)
	}
}

func doMain() error {

	_, err := config.Initialize()
	if err != nil {
		return err
	}

	ctx := context.Background()
	client, err := session.NewSession(ctx)
	if err != nil {
		return err
	}

	out, err := bsky.FeedGetTimeline(ctx, client, "reverse-chronological", "", 10)
	if err != nil {
		return fmt.Errorf("get timeline failed: %w", err)
	}

	for _, item := range out.Feed {
		if item.Post == nil {
			continue
		}
		fmt.Println(format(item))
	}

	return nil
}

func format(feed *bsky.FeedDefs_FeedViewPost) string {

	// Record.Val は CBORMarshaler 型なので、*bsky.FeedPost にキャスト
	txt := ""
	if rec, ok := feed.Post.Record.Val.(*bsky.FeedPost); ok {
		txt = *feed.Post.Author.DisplayName + ":" + rec.Text
	}

	return txt
}
