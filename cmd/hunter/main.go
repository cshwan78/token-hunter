package main

import (
	"fmt"
	"log"

	"github.com/cshwan78/pkg/api/github"
)

func main() {
	fmt.Printf("%s", "hello world")
	res, ok := github.APISearch("xoxb-")
	if ok != nil {
		log.Fatalln(ok)
		return
	}
	fmt.Printf("%s", res)
}
