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
	client := session.NewSession(ctx)

	out, err := bsky.FeedGetTimeline(ctx, client, "reverse-chronological", "", 10)
	if err != nil {
		return fmt.Errorf("get timeline failed: %w", err)
	}

	for _, item := range out.Feed {
		if item.Post == nil {
			continue
		}
		// Record.Val は CBORMarshaler 型なので、*bsky.FeedPost にキャスト
		if rec, ok := item.Post.Record.Val.(*bsky.FeedPost); ok {
			fmt.Println(*item.Post.Author.DisplayName, ":", rec.Text)
		}
	}

	return nil
}
