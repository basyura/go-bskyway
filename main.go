package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"bskyway/config"
	"bskyway/model"
	"bskyway/notifier"
	"bskyway/session"

	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/common-nighthawk/go-figure"
)

func main() {
	if err := doMain(); err != nil {
		fmt.Println(err)
	}
}

func doMain() error {
	// 設定の初期化
	if err := initialize(); err != nil {
		return err
	}

	// セッションの生成
	ctx := context.Background()
	client, err := session.NewSession(ctx)
	if err != nil {
		return err
	}

	// 最後に通知したフィード
	lastCid := ""
	for {
		// タイムラインの取得
		out, err := bsky.FeedGetTimeline(ctx, client, "reverse-chronological", "", 10)
		// セッションの再生成
		if err != nil {
			if strings.Contains(err.Error(), "ExpiredToken") {
				client, err = session.NewSession(ctx)
				if err != nil {
					return fmt.Errorf("session refresh failed: %w", err)
				}
				continue
			}
		}

		// 通知
		lastCid = notify(out, lastCid)
		// 次の取り込みまでの待機
		time.Sleep(30 * time.Second)
	}
}

func initialize() error {
	fmt.Println("")
	myFigure := figure.NewFigure("bskyway", "", true)
	myFigure.Print()
	// 設定の初期化
	config, err := config.Initialize()
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("identifier   :", config.Identifier)
	fmt.Println("IconCacheDir :", config.IconCacheDir)
	fmt.Println("")

	return nil
}

func notify(out *bsky.FeedGetTimeline_Output, lastCid string) string {
	// 最後に通知したフィード以降のみに絞る
	feeds := filter(out, lastCid)
	// 通知対象が無い
	if len(feeds) == 0 {
		return lastCid
	}

	// 取得と逆順に古いものから通知
	for i := len(feeds) - 1; i >= 0; i-- {
		feed := feeds[i]
		if feed.Post == nil {
			continue
		}

		post := model.ConvertToPost(feed)
		lastCid = post.Cid

		// 通知
		notifier.Notify(post)

		fmt.Println(post.Format())
		fmt.Println("--------------------------------------------------")
		time.Sleep(5 * time.Second)
	}

	return lastCid
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
