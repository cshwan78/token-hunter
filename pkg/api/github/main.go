package github

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
)

type searchResult struct {
	repoName      string
	filePath      string
	githubWebPath string
	snippet       string
}

func GetWebSearchResult(keyword string) (string, error) {
	req, err := http.NewRequest("GET", "https://github.com/search?o=desc&s=indexed&type=Code", nil)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	//req.Header.Set("Accept", "application/vnd.github.v3+json")/

	q := req.URL.Query()
	q.Add("q", keyword)
	req.URL.RawQuery = q.Encode()

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	return string(body), nil
}

func ParseSearchHTML(htmlbody string) ([]searchResult, error) {
	doc, err := htmlquery.Parse(strings.NewReader(htmlbody))

	if err != nil {
		return nil, err
	}

	var repoListXPath string = "//*[@id='code_search_results']/div[1]/div/div[1]"
	repoList := htmlquery.Find(doc, repoListXPath)
	var searchResults []searchResult = make([]searchResult, 0)

	for _, node := range repoList {
		var snippet string = ""

		repoNameNode := htmlquery.FindOne(node, "/div[1]/a")
		if repoNameNode == nil {
			continue
		}
		repoName := htmlquery.SelectAttr(repoNameNode, "href")

		filePathNode := htmlquery.FindOne(node, "/div[2]/a")
		if filePathNode == nil {
			continue
		}
		filePath := htmlquery.SelectAttr(filePathNode, "title")
		githubWebPath := htmlquery.SelectAttr(filePathNode, "href")

		snippetNode := htmlquery.FindOne(node, "/div[3]")
		if snippetNode == nil {
			continue
		}
		snippetLinesNodes := htmlquery.Find(node, "/div[3]//td[@class='blob-code blob-code-inner']")
		for _, lineNode := range snippetLinesNodes {
			snippet += htmlquery.InnerText(lineNode)
		}
		searchResult := searchResult{repoName, filePath, githubWebPath, snippet}
		searchResults = append(searchResults, searchResult)
	}

	return searchResults, nil
}

// func GetSearchResult(keyword string) error {
// 	htmlbody, err := GetWebSearchResult(keyword)
// 	if err != nil {
// 		return nil
// 	}
//}
