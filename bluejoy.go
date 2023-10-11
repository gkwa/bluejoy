package bluejoy

import (
	"encoding/gob"
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"syscall"
	"time"

	"github.com/adrg/xdg"
	gocache "github.com/patrickmn/go-cache"
)

var configRelPath = "bluejoy/keys.db"

func Main() int {
	cache := gocache.New(1*time.Minute, 2*time.Minute)

	path, _ := genCachePath(configRelPath)
	slog.Debug("cache", "path", path)

	checkFileExists(path)
	os.Remove(path)
	checkFileExists(path)

	cacheItem := PushbulletHTTReply{
		Pushes: []Push{
			{URL: "https://news.ycombinator.com/"},
			{URL: "https://go.dev/blog/gob"},
		},
	}
	cache.Set("foo", cacheItem, gocache.DefaultExpiration)

	// gob stuff to save cache:
	cacheSnapshot := cache.Items()
	file2, _ := os.Create(path)
	defer file2.Close()
	encoder := gob.NewEncoder(file2)
	encoder.Encode(cacheSnapshot)

	// ensure cache has been saved:
	checkFileExists(path)

	// pretend to restart app and load cache from file3:
	file3, err := os.Open(path)
	if err != nil {
		slog.Debug("file access", "error", err.Error())
		return 1
	}
	defer file3.Close()

	var newCache2 gocache.Cache
	decoder := gob.NewDecoder(file3)

	if err := decoder.Decode(&newCache2); err != nil {
		slog.Debug("decode", "error", err.Error())
		return 1
	}

	slog.Debug("cache", "itemCount", newCache2.ItemCount())

	r, future, found := newCache2.GetWithExpiration("foo")

	if found {
		slog.Debug("mymessage", "result", r.(PushbulletHTTReply))
		slog.Debug("duration", "n", time.Until(future).Truncate(time.Second))
	}

	foo, found := newCache2.Get("foo")
	if found {
		slog.Debug("debug", "foo", foo.(PushbulletHTTReply))
	}
	return 0
}

func persistReply(reply PushbulletHTTReply, path string) error {
	file, _ := os.Create(path)
	defer file.Close()
	encoder := gob.NewEncoder(file)
	encoder.Encode(reply)

	return nil
}

func loadCache(path string) error {
	file, err := os.Open(path)
	if err != nil {
		slog.Debug("file access", "error", err.Error())
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	var reply PushbulletHTTReply

	if err := decoder.Decode(&reply); err != nil {
		slog.Debug("decode", "error", err.Error())
		return err
	}

	slog.Debug("user", "email", reply.Pushes[0].SenderEmail)

	return nil
}

func genCachePath(configRelPath string) (string, error) {
	configFilePath, err := xdg.ConfigFile(configRelPath)
	if err != nil {
		return "", err
	}

	dirPerm := os.FileMode(0o700)

	d := filepath.Dir(configFilePath)

	if err := os.MkdirAll(d, dirPerm); err != nil {
		slog.Error("cache", "mkdir", "error", err.Error())
		return "", err
	}

	slog.Debug("cache", "path", configFilePath)
	return configFilePath, nil
}

func logPathStats(filePath string) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		slog.Error("stat", "path", filePath, "error", err.Error())
		return
	}

	fileUID := fileInfo.Sys().(*syscall.Stat_t).Uid

	// Use the user package to get the user information
	u, err := user.LookupId(fmt.Sprintf("%d", fileUID))
	if err != nil {
		slog.Error("user info", "user", u, "error", err.Error())
		return
	}

	slog.Debug("owner", "path", filePath, "user", u.Username)
}
