package main

import (
	"context"
	"fmt"
	"time"

	"bskyway/config"
	"bskyway/model"
	"bskyway/session"

	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/go-toast/toast"
)

func main() {
	if err := doMain(); err != nil {
		fmt.Println(err)
	}
}

func doMain() error {
	// 設定の初期化
	if _, err := config.Initialize(); err != nil {
		return err
	}

	// セッションの生成
	ctx := context.Background()
	client, err := session.NewSession(ctx)
	if err != nil {
		return err
	}

	lastCid := ""
	for {
		out, err := bsky.FeedGetTimeline(ctx, client, "reverse-chronological", "", 5)
		if err != nil {
			return fmt.Errorf("get timeline failed: %w", err)
		}

		feeds := filter(out, lastCid)
		if len(feeds) != 0 {
			for i := len(feeds) - 1; i >= 0; i-- {
				feed := feeds[i]
				if feed.Post == nil {
					continue
				}

				post := model.ConvertToPost(feed)
				if post.Cid == lastCid {
					break
				}
				lastCid = post.Cid

				fmt.Println(post.Format())
				fmt.Println("--------------------------------------------------")
				time.Sleep(5 * time.Second)

				notification := toast.Notification{
					AppID:   "bskyway",
					Title:   post.Name,
					Message: post.Text,
				}
				notification.Push()
			}
		}
		// fmt.Println("==================================================")

		// 次の取り込み用
		time.Sleep(30 * time.Second)
	}
}

func filter(out *bsky.FeedGetTimeline_Output, lastCid string) []*bsky.FeedDefs_FeedViewPost {
	feeds := []*bsky.FeedDefs_FeedViewPost{}
	for _, feed := range out.Feed {
		if feed.Post.Cid == lastCid {
			break
		}
		feeds = append(feeds, feed)
	}

	return feeds
}
