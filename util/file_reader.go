package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("/Users/manishjain/go-redis/src/cache/example/data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		words := strings.Split(s, "and")
		if len(words) != 2 {
			fmt.Printf("error in splitting")
			return
		}
		first := strings.Trim(words[0], " ")
		second := strings.Trim(words[1], " ")
		fmt.Printf("%s:%s\n", first, second)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
