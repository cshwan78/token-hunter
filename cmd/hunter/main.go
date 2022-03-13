package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/cshwan78/pkg/api/github"
)

func main() {
	res, ok := github.GetSearchResult("xoxb-")
	if ok != nil {
		log.Fatalln(ok)
		return
	}
	// fmt.Printf("%+v", res)

	file, err := os.Create("./output.csv")
	if err != nil {
		panic(err)
	}
	wr := csv.NewWriter(file)

	for _, searchResult := range res {
		wr.Write([]string{searchResult.RepoName, searchResult.FilePath, searchResult.GithubWebPath}) //, searchResult.Snippet})
	}

	wr.Flush()
}
