//go:build darwin

package notifier

import (
	"bskyway/config"
	"bskyway/model"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func Notify(post *model.Post) {
	title := post.Name
	message := post.Text

	avatarURL := post.Avatar
	// Generate a deterministic filename based on the avatar URL so we can reuse it from the temp directory.
	hash := sha1.Sum([]byte(avatarURL))
	filename := fmt.Sprintf("%x%s", hash[:], ".jpeg")
	cacheDir := config.Instance().IconCacheDir
	cachePath := filepath.Join(cacheDir, filename)

	// Download the avatar only if it does not already exist in the temp directory.
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		resp, err := http.Get(avatarURL)
		if err != nil {
			fmt.Println("avatar download error:", err)
		} else {
			defer resp.Body.Close()
			out, err := os.Create(cachePath)
			if err != nil {
				fmt.Println("avatar file create error:", err)
			} else {
				if _, err = io.Copy(out, resp.Body); err != nil {
					fmt.Println("avatar save error:", err)
				}
				out.Close()
			}
		}
	}

	fmt.Println("cachePath :", cachePath)
	cmd := exec.Command("terminal-notifier",
		"-title", title,
		"-message", message,
		"-appIcon", cachePath,
		"-contentImage", cachePath,
	)

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}
