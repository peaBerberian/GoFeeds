// This file describes the feedCache type used to store the RSS fetched
// from various RSS URLs
// A feedCache have 4 methods:
//   - has:   Returns true if a cache is present for a specific Id.
//   - get:   Returns the cache for the corresponding Id.
//   - set:   Set the cache for the corresponding Id.
//   - reset: Re-initalize the entire cache.

package main

import "fmt"
import "time"
import "sync"
import "errors"

const MAX_CHACHE_ELEMENTS = 100

type feedCache struct {
	websites []websiteCache
	mutex    sync.Mutex
}

type websiteCache struct {
	Id    int        // linked to the ids in config.json
	Date  time.Time  // last update date
	Cache feedFormat // What was set at this date
}

func (w *feedCache) has(id int) bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	for _, webs := range w.websites {
		if webs.Id == id {
			return true
		}
	}
	return false
}

func (w *feedCache) get(id int) (websiteCache, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	for _, web := range w.websites {
		if web.Id == id {
			return web, nil
		}
	}
	var errorText = fmt.Sprintf("Could not find cache for the Id %d.", id)
	return websiteCache{}, errors.New(errorText)
}

func (w *feedCache) set(id int, data feedFormat) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	for _, webs := range w.websites {
		if webs.Id == id {
			webs.Cache = data
			webs.Date = time.Now()
			return nil
		}
	}
	newElem := websiteCache{Id: id,
		Cache: data, Date: time.Now()}
	if len(w.websites) >= MAX_CHACHE_ELEMENTS {
		w.websites[0] = newElem
		return nil
	}
	w.websites = append(w.websites, newElem)
	return nil
}

// Re-init the cache
func (w *feedCache) reset() {
	w.mutex.Lock()
	w.websites = []websiteCache{}
	w.mutex.Unlock()
}
