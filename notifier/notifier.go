package notifier

import (
	"bskyway/config"
	"bskyway/model"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func downloadAvatar(post *model.Post) string {
	avatarURL := post.Avatar
	// Generate a deterministic filename based on the avatar URL so we can reuse it from the temp directory.
	hash := sha1.Sum([]byte(avatarURL))
	filename := fmt.Sprintf("%x%s", hash[:], ".jpeg")
	cacheDir := config.Instance().IconCacheDir
	imgPath := filepath.Join(cacheDir, filename)

	// Download the avatar only if it does not already exist in the temp directory.
	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		resp, err := http.Get(avatarURL)
		if err != nil {
			fmt.Println("avatar download error:", err)
		} else {
			defer resp.Body.Close()
			out, err := os.Create(imgPath)
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

	return imgPath
}
