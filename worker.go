package cache

import (
	"fmt"
	"io"
	"net/http"
)

type Http_Req struct {
	W   *http.ResponseWriter
	Req *http.Request
}

func Http_worker(work chan *Http_Req, i int) {
	for {
		work_item, ok := <-work
		if !ok {
			if debug {
				fmt.Println("chanel %d closed", i)
			}
			break
		}
		if debug {
			fmt.Printf("%s recvd on %d\n", work_item.Req.URL.Path[1:], i)
		}
		key := work_item.Req.URL.Path[1:]
		if debug {
			fmt.Println("query key:%s", key)
		}
		value, bfailed := Handle_get_request(key)
		if bfailed == true {
			if debug {
				fmt.Println(" error for key:%s", key)
				io.WriteString(*work_item.W, "error")
			}
		} else {
			if debug {
				fmt.Println("%s:%s", key, value)
				fmt.Fprintf(*work_item.W, "%s is love of %s", key, value)
				io.WriteString(*work_item.W, value)
			}
		}
	}

	return
}
