package bluejoy

import (
	"log/slog"
	"os"
)

func checkFileExists(path string) {
	if _, err := os.Stat(path); err == nil {
		slog.Debug("checking cache file", "path", path, "exists", true)
	} else if os.IsNotExist(err) {
		slog.Debug("checking cache file", "path", path, "exists", false)
	} else {
		slog.Error("checking for file", "path", path, "error", err.Error())
	}
}
