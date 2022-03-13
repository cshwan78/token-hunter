package github

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
)

type searchResult struct {
	RepoName      string
	FilePath      string
	GithubWebPath string
	Snippet       string
}

func GetWebSearchResult(keyword string, page int) (string, error) {
	req, err := http.NewRequest("GET", "https://github.com/search?o=desc&s=indexed&type=Code", nil)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	userSession, ok := os.LookupEnv("GITHUB_USER_SESSION")
	if !ok {
		return "", errors.New("no github token found")
	}

	req.AddCookie(&http.Cookie{Name: "user_session", Value: userSession})
	//req.Header.Set("Accept", "application/vnd.github.v3+json")/

	q := req.URL.Query()
	q.Add("q", keyword)
	q.Add("p", fmt.Sprint(page))

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

func GetSearchResult(keyword string) ([]searchResult, error) {
	htmlbody, err := GetWebSearchResult(keyword, 0)
	if err != nil {
		return nil, err
	}
	doc, err := htmlquery.Parse(strings.NewReader(htmlbody))
	searchResultCountNode := htmlquery.FindOne(doc, "//div[@class='col-12 col-md-9 float-left px-2 pt-3 pt-md-0 codesearch-results']/div/div/h3")

	if searchResultCountNode == nil {
		return nil, errors.New("No searchResultCountNode found")
	}
	searchResultCountString := htmlquery.InnerText(searchResultCountNode)

	searchResultCountString = strings.TrimLeft(searchResultCountString, " \n\t")
	searchResultCountString = strings.Split(searchResultCountString, " code ")[0]
	searchResultCountString = strings.Replace(searchResultCountString, ",", "", -1)
	searchResultCount, err := strconv.Atoi(searchResultCountString)

	if err != nil {
		return nil, err
	}

	totalPage := int(math.Min(5, float64(searchResultCount)/10.0))
	var totalSearchResult []searchResult
	searchResults, err := ParseSearchHTML(htmlbody)

	if err != nil {
		return nil, err
	}

	totalSearchResult = append(totalSearchResult, searchResults...)

	for i := 1; i < totalPage; i += 1 {
		htmlbody, err := GetWebSearchResult(keyword, i)
		if err != nil {
			break
		}
		searchResults, err := ParseSearchHTML(htmlbody)
		if err != nil {
			break
		}
		totalSearchResult = append(totalSearchResult, searchResults...)

	}
	return totalSearchResult, nil
}
