package cache

import (
	"fmt"
	"io"
	"net/http"
)

func Http_worker(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Path[1:]
	if debug {
		fmt.Println("query key:%s", key)
	}
	value, bfailed := Handle_get_request(key)
	if bfailed == true {
		if debug {
			fmt.Println(" error for key:%s", key)
		}
		io.WriteString(w, "error")
	} else {
		if debug {
			fmt.Println("%s:%s", key, value)
		}
		io.WriteString(w, value)
	}

	return
}
