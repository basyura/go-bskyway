//go:build darwin

package notifier

import (
	"bskyway/model"
	"fmt"
	"os/exec"
)

func Notify(post *model.Post) {
	title := post.Name
	message := post.Text
	imgPath := downloadAvatar(post)

	cmd := exec.Command("terminal-notifier",
		"-title", title,
		"-message", message,
		"-contentImage", imgPath,
	)

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}
