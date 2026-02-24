package core

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("serveur a répondu : %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func removeFile(filepath string) error {
	if err := os.Remove(filepath); err != nil {
		return err
	}
	return nil
}

func getVersionFromPath(path string) string {
	base := filepath.Base(path)
	base, _ = strings.CutSuffix(base, ".exe")
	idx := strings.LastIndex(base, "_")
	if idx == -1 {
		return "unknown"
	}
	return base[idx+1:]
}
