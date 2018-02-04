package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "http://localhost:8080/"
	relations := map[string]string{
		"Abbott":  "Costello",
		"Adam":    "Eve",
		"apples":  "oranges",
		"Antony":  "Cleopatra",
		"bacon":   "eggs",
		"Barbie":  "Ken",
		"Batman":  "Robin",
		"bed":     "breakfast",
		"before":  "after",
		"Bert":    "Ernie",
		"big":     "small",
		"black":   "white",
		"Bonnie":  "Clyde",
		"bow":     "arrow",
		"boys":    "girls",
		"bread":   "butter",
		"Cain":    "Abel",
		"Castor":  "Pollux",
		"cold":    "hot",
		"Crick":   "Watson",
		"cut":     "paste",
		"day":     "night",
		"death":   "taxes",
		"Dick":    "Jane",
		"divide":  "conquer",
		"dogs":    "cats",
		"each":    "every",
		"ebony":   "ivory",
		"fast":    "slow",
		"fat":     "thin",
		"fife":    "drum",
		"fire":    "brimstone",
		"fish":    "chips",
		"flotsam": "jetsam",
		"free":    "clear",
	}
	for key, val := range relations {
		req := url + key
		resp, err := http.Get(req)
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if string(body) == val {
			fmt.Printf("PASSED\n")
		} else {
			fmt.Printf("FAILED.%s:%s (%s)\n", key, val, string(body))
		}
		resp.Body.Close()
	}

}
