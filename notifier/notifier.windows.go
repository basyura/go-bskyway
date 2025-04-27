//go:build windows

package notifier

import (
	"bskyway/model"

	"github.com/go-toast/toast"
)

func Notify(post *model.Post) {

	notification := toast.Notification{
		AppID:   "bskyway",
		Title:   post.Name,
		Message: post.Text,
	}
	notification.Push()
}
