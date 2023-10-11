package bluejoy

import (
	"log/slog"
	"os"
)

func checkFileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if os.IsNotExist(err) {
		slog.Debug("checking cache file", "path", path, "exists", false)
		return false
	} else {
		slog.Error("checking for file", "path", path, "error", err.Error())
		panic(err)
	}
}
