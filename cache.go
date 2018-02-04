package cache

import (
	"container/list"
	"fmt"
	"log"
	"redis"
	"sync"
	"time"
)

type node struct {
	key, value string
}
type read_cache struct {
	capacity    int
	key_storage map[string]*list.Element
	list        *list.List
	backend     redis.Client
	timer       *time.Ticker
	lock        sync.Mutex
	secs        uint16
}

var debug bool = false
var cache read_cache

func delete_item_from_cache() {
	/*remove the least used key from the cache*/
	pEl := cache.list.Back()
	item := cache.list.Remove(pEl).(node)
	if debug {
		fmt.Println("removing key:%s from cache", item.key)
	}
	delete(cache.key_storage, item.key)
}

func process_timer_expired() {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	if debug {
		fmt.Println("inside process_timer_expired.len:%d", cache.list.Len())
	}
	len := cache.list.Len()
	if len > 1 {
		delete_item_from_cache()
	}
	return
}
func start_timer() {
	for _ = range cache.timer.C {
		process_timer_expired()
	}
}
func request_key_from_redis(key string) (item node, bfailed bool) {
	if debug {
		fmt.Println("requesting key %s from redis", key)
	}
	bfailed = false
	value, err := cache.backend.Get(key)
	if err != nil {
		bfailed = true
		if debug {
			log.Println("error on get at redis", err)
		}
		return
	}

	if value == nil {
		if debug {
			fmt.Println("no key at backend")
		}
		bfailed = true
		return
	}
	item.key = key
	item.value = fmt.Sprintf("%s", value)
	return
}

func Handle_get_request(key string) (str string, bFailed bool) {
	/*first check in the cache whether*/
	var item node
	cache.lock.Lock()
	defer cache.lock.Unlock()
	value, err := cache.key_storage[key]
	if err == false || value == nil {
		if debug {
			fmt.Println("key %s not present in cache", key)
		}
		item, bFailed = request_key_from_redis(key)
		if bFailed == true {
			return
		}

	} else {
		if debug {
			fmt.Println("key %s  present in cache", key)
		}
		/*Removing from list and key storage to again rehash*/
		item = cache.list.Remove(value).(node)
		delete(cache.key_storage, item.key)
	}
	if debug {
		fmt.Println("adding item key:%s val:%s to cache", item.key, item.value)
	}

	if cache.list.Len()+1 > cache.capacity {
		if debug {
			fmt.Println("cache full len:%d capacity:%d", cache.list.Len(), cache.capacity)
		}
		delete_item_from_cache()
	}
	/*Making it most recently used  item*/
	ele := cache.list.PushFront(item)
	cache.key_storage[key] = ele
	str = item.value
	return
}

func Init_cache(cli redis.Client, max_keys int, expirtytime uint16) {
	cache.backend = cli
	cache.capacity = 128
	cache.key_storage = make(map[string]*list.Element)
	cache.list = list.New()
	cache.secs = expirtytime
	cache.timer = time.NewTicker(time.Second * time.Duration(cache.secs))
	go start_timer()
	return
}
