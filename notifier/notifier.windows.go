//go:build windows

package notifier

import (
	"bskyway/model"

	"github.com/gen2brain/beeep"
)

func Notify(post *model.Post) {
	path := downloadAvatar(post)
	beeep.Notify(post.Name, post.Text, path)
}
