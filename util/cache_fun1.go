package main

import (
	"cache"
	"flag"
	"fmt"
	"log"
	"net/http"
	"redis"
	"runtime"
)

func main() {

	// Parse command-line flags; needed to let flags used by Go-Redis be parsed.
	flag.Parse()

	// create the client.  Here we are using a synchronous client.
	// Using the default ConnectionSpec, we are specifying the client to connect
	// to db 13 (e.g. SELECT 13), and a password of go-redis (e.g. AUTH go-redis)

	spec := redis.DefaultSpec().Db(13).Password("")
	client, e := redis.NewSynchClientWithSpec(spec)
	if e != nil {
		log.Println("failed to create the client", e)
		return
	}
	items := map[string]string{
		"manish":  "ankita",
		"abhinav": "ritika",
		"cisco":   "arachna",
		"satish":  "varsha"}

	for key, value := range items {
		fmt.Println("storing key:%s,val:%s", key, value)
		bytesval := []byte(value)
		client.Set(key, bytesval)
	}
	cache.Init_cache(client)
	valString, bfailed := cache.Handle_get_request("manish")
	if bfailed == true {
		fmt.Println("error on while get :%s", "manish")
		return
	} else {
		fmt.Println("got key:%s val:%s", "manish", valString)
	}

	valString, bfailed = cache.Handle_get_request("abhinav")
	if bfailed == true {
		fmt.Println("error on while get :%s", "abhinav")
		return
	} else {
		fmt.Println("got key:%s val:%s", "abhinav", valString)
	}

	valString, bfailed = cache.Handle_get_request("cisco")
	if bfailed == true {
		fmt.Println("error on while get :%s", "cisco")
		return
	} else {
		fmt.Println("got key:%s val:%s", "cisco", valString)
	}

	fmt.Println("max procs:%d", runtime.GOMAXPROCS(-1))
	//	setup_worker_buff()
	http.HandleFunc("/", cache.Http_worker)
	http.ListenAndServe(":8080", nil)
}
