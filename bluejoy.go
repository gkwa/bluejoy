package bluejoy

import (
	"encoding/gob"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	gocache "github.com/patrickmn/go-cache"
)

var configRelPath = "bluejoy/keys.db"

func Main() int {
	path, _ := genCachePath(configRelPath)
	slog.Debug("cache", "path", path)

	cache1 := gocache.New(3*time.Minute, 4*time.Minute)
	slog.Debug("cache", "exists", checkFileExists(path))

	slog.Debug("cache", "action", "deleting cache file")

	// ensure we're starting clean:
	os.Remove(path)
	slog.Debug("cache", "exists", checkFileExists(path))

	cacheItem := PushbulletHTTReply{
		Pushes: []Push{
			{URL: "https://news.ycombinator.com/"},
			{URL: "https://go.dev/blog/gob"},
		},
	}
	cache1.Set("foo", cacheItem, 2*time.Minute)
	slog.Debug("check in memory cache items", "count", cache1.ItemCount())

	cacheSnapshot := cache1.Items()

	gob.Register(PushbulletHTTReply{})

	// serialize using gob:
	file, _ := os.Create(path)
	encoder := gob.NewEncoder(file)
	err := encoder.Encode(cacheSnapshot)
	if err != nil {
		slog.Error("encode", "error", err.Error())
	}
	defer file.Close()
	slog.Debug("checking existance of file cache", "exists", checkFileExists(path))

	// unmarshal cache from file
	file2, err := os.Open(path)
	if err != nil {
		slog.Debug("file access", "error", err.Error())
		return 1
	}
	defer file2.Close()

	decoder := gob.NewDecoder(file2)

	var cacheMap map[string]gocache.Item
	if err := decoder.Decode(&cacheMap); err != nil {
		slog.Debug("decode", "error", err.Error())
		return 1
	}

	cache2 := gocache.NewFrom(1*time.Minute, 2*time.Minute, cacheMap)
	reply, future, found := cache2.GetWithExpiration("foo")

	expires := time.Until(future).Truncate(time.Second)
	e := reply.(PushbulletHTTReply)
	slog.Debug("newCache", "found", found, "expires", expires, "now", time.Now(), "future", future, "reply", e.Pushes[0].URL)

	return 0
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

	return configFilePath, nil
}
