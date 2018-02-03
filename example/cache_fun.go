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

//var channels []chan *cache.Http_Req
var buff_chan chan cache.Http_Req

//var cases []reflect.SelectCase
//var incoming_req cache.Http_Req

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
	setup_worker_buff()
	http.HandleFunc("/", handler_buff)
	http.ListenAndServe(":8080", nil)
}

/*
func shut_down_worker() {
	for i := 0; i < len(channels); i++ {
		close(channels[i])
	}
}*/
func setup_worker_buff() {
	/*select case,array of channel and workser*/
	n := runtime.GOMAXPROCS(-1)
	buff_chan = make(chan cache.Http_Req, n)
	for i := 0; i < n; i++ {
		go cache.Http_worker(buff_chan, i)
	}
}

/*
func setup_worker() {
	n := runtime.GOMAXPROCS(-1)
	channels = make([]chan *cache.Http_Req, n)
	cases = make([]reflect.SelectCase, n)
	for i := 0; i < n; i++ {
		channels[i] = make(chan *cache.Http_Req)
		cases[i] = reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(channels[i]), Send: reflect.ValueOf(&incoming_req)}
		go cache.Http_worker(channels[i], i)
	}
}*/

/*
func handler(w http.ResponseWriter, r *http.Request) {
	incoming_req = cache.Http_Req{W: &w, Req: r}
	fmt.Fprintf(*incoming_req.W, "Hi there, I love %s!", r.URL.Path[1:])
	ch, _, _ := reflect.Select(cases)
	fmt.Printf("channel :%d is chooses", ch)
}*/
func handler_buff(w http.ResponseWriter, r *http.Request) {
	incoming_req := cache.Http_Req{W: w, Req: r}
	//	fmt.Fprintf(incoming_req.W, "Hi there, I love %s!", r.URL.Path[1:])
	buff_chan <- incoming_req
}
