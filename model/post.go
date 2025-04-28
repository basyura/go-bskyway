package model

import (
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api/bsky"
)

type Post struct {
	Name      string
	Text      string
	Cid       string
	CreatedAt time.Time
	Avatar    string
}

func (p *Post) Format() string {

	buf := ""
	for i, line := range strings.Split(p.Text, "\n") {
		if i != 0 {
			buf += "\n"
		}

		// buf += runewidth.FillLeft("", 4) + txt
		buf += line
	}

	name := p.Name
	// name = runewidth.Truncate(name, 20, "")
	// name = runewidth.FillRight(name, 20)

	return name + " - " + p.CreatedAt.Format("2006/01/02 15:04:05") + "\n" + buf
}

func ConvertToPost(feed *bsky.FeedDefs_FeedViewPost) *Post {

	name := *feed.Post.Author.DisplayName

	post := &Post{
		Name:   name,
		Cid:    feed.Post.Cid,
		Avatar: *feed.Post.Author.Avatar,
	}

	// Record.Val は CBORMarshaler 型なので、*bsky.FeedPost にキャスト
	if rec, ok := feed.Post.Record.Val.(*bsky.FeedPost); ok {
		post.Text = rec.Text
		post.CreatedAt = convertTime(rec.CreatedAt)
	}

	return post
}

func convertTime(s string) time.Time {

	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		panic(err)
	}

	// システム設定からロケーションを自動取得
	loc := time.Now().Location()
	local := t.In(loc)

	return local
}
