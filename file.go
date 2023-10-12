package bluejoy

import (
	"log/slog"
	"os"
)

func checkFileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if os.IsNotExist(err) {
		slog.Debug("checking file existence", "path", path, "exists", false)
		return false
	} else {
		slog.Error("checking file type", "path", path, "error", err.Error())
		panic(err)
	}
}
