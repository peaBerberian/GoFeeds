package requests

import "fmt"
import "time"
import "sync"
import "errors"
import "github.com/peaberberian/OscarGoGo/format"

const MAX_CHACHE_ELEMENTS = 100

type feedCache struct {
	websites []websiteCache
	timeout  int
	mutex    sync.Mutex
}

type websiteCache struct {
	id    int               // linked to the ids in config.json
	date  time.Time         // last update date
	cache format.FeedFormat // What was set at this date
}

func isDeprecated(date time.Time, timeout int) bool {
	var now = time.Now()
	return timeout*1000 < now.Nanosecond()-date.Nanosecond()
}

func (w *feedCache) cleanDeprecated() {
	w.mutex.Lock()
	var newWebsiteCache []websiteCache
	for i, webs := range w.websites {
		if isDeprecated(webs.date, w.timeout) {
			// delete from slice
			newWebsiteCache = append(w.websites[:i], w.websites[i+1:]...)
		}
	}
	w.websites = newWebsiteCache
	w.mutex.Unlock()
}

func NewCache(timeout int) feedCache {
	return feedCache{timeout: timeout}
}

func (w *feedCache) HasCacheForId(id int) bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	for _, webs := range w.websites {
		if webs.id == id {
			if isDeprecated(webs.date, w.timeout) {
				return false
			}
			go w.cleanDeprecated()
			return true
		}
	}
	return false
}

func (w *feedCache) GetCacheForId(id int) (format.FeedFormat, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	for _, web := range w.websites {
		if web.id == id {
			if !isDeprecated(web.date, w.timeout) {
				return web.cache, nil
			} else {
				go w.cleanDeprecated()
			}
		}
	}
	var errorText = fmt.Sprintf("Could not find cache for the Id %d.", id)
	return format.FeedFormat{}, errors.New(errorText)
}

func (w *feedCache) SetCacheForId(id int, data format.FeedFormat) error {
	w.mutex.Lock()

	// will execute cleanDeprecated after unlocking the mutex
	go w.cleanDeprecated()
	defer w.mutex.Unlock()
	for _, webs := range w.websites {
		if webs.id == id {
			webs.cache = data
			webs.date = time.Now()
			return nil
		}
	}
	newElem := websiteCache{id: id,
		cache: data, date: time.Now()}
	if len(w.websites) >= MAX_CHACHE_ELEMENTS {
		w.websites[0] = newElem
		return nil
	}
	w.websites = append(w.websites, newElem)
	return nil
}

// Re-init the cache
func (w *feedCache) ResetCache() {
	w.mutex.Lock()
	w.websites = []websiteCache{}
	w.mutex.Unlock()
}
