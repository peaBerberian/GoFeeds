// This file describes the feedCache type used to store the RSS fetched
// from various RSS URLs
// A feedCache have 4 methods:
//   - has:   Returns true if a cache is present for a specific id.
//   - get:   Returns the cache for the corresponding id.
//   - set:   Set the cache for the corresponding id.
//   - reset: Re-initalize the entire cache.

package requests

import "fmt"
import "time"
import "sync"
import "errors"
import "github.com/peaberberian/OscarGoGo/format"

const MAX_CHACHE_ELEMENTS = 100

type FeedCache struct {
	websites []websiteCache
	mutex    sync.Mutex
}

type websiteCache struct {
	id    int               // linked to the ids in config.json
	date  time.Time         // last update date
	cache format.FeedFormat // What was set at this date
}

func (w *FeedCache) has(id int) bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	for _, webs := range w.websites {
		if webs.id == id {
			return true
		}
	}
	return false
}

func (w *FeedCache) get(id int) (websiteCache, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	for _, web := range w.websites {
		if web.id == id {
			return web, nil
		}
	}
	var errorText = fmt.Sprintf("Could not find cache for the Id %d.", id)
	return websiteCache{}, errors.New(errorText)
}

func (w *FeedCache) set(id int, data format.FeedFormat) error {
	w.mutex.Lock()
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
func (w *FeedCache) reset() {
	w.mutex.Lock()
	w.websites = []websiteCache{}
	w.mutex.Unlock()
}
