package main

import (
	"bufio"
	"cache"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"redis"
	"strconv"
	"strings"
)

const (
	datafile   string = "data.txt"
	configfile string = "proxy.conf"
)

type config struct {
	capacity int
	secs     uint16
	port     string
}

func read_config() (conf config, bfailed bool) {
	file, err := os.Open(configfile)
	bfailed = false
	if err != nil {
		log.Fatal(err)
		bfailed = true
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if s[0] == '#' {
			continue
		}
		words := strings.Split(s, ":")
		if len(words) != 2 {
			fmt.Printf("error in splitting:%s\n", s)
			bfailed = true
			return
		}
		first := strings.Trim(words[0], " ")
		second := strings.Trim(words[1], " ")
		number := 0
		if first == "maxkeys" || first == "expirytime" {
			number, err = strconv.Atoi(second)
			if err != nil {
				bfailed = true
				return
			}
		}
		switch first {
		case "maxkeys":
			conf.capacity = number
		case "expirytime":
			conf.secs = uint16(number)
		case "port":
			conf.port = second
		default:
			fmt.Printf("not valid.%s\n", first)
			bfailed = true
		}
		if bfailed {
			return
		}
	}
	//	fmt.Printf("cap:%d,sec:%d,port:%s\n", conf.capacity, conf.secs, conf.port)
	return
}
func setup_backend_data(client redis.Client) (bfailed bool) {
	file, err := os.Open(datafile)
	bfailed = false
	if err != nil {
		log.Fatal(err)
		bfailed = true
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		words := strings.Split(s, ":")
		if len(words) != 2 {
			fmt.Printf("error in splitting")
			bfailed = true
			return
		}
		first := strings.Trim(words[0], " ")
		second := strings.Trim(words[1], " ")
		//fmt.Printf("%s,%s", first, second)
		client.Set(first, []byte(second))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		bfailed = true
	}
	return
}
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
	setup_backend_data(client)
	conf, bfailed := read_config()
	if bfailed {
		fmt.Printf("error while reading config.aborting\n")
		return
	}
	cache.Init_cache(client, conf.capacity, conf.secs)
	http.HandleFunc("/", cache.Http_worker)
	address := ":" + conf.port
	http.ListenAndServe(address, nil)
}
