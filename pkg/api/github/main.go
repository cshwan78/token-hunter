package github

import (
	"io/ioutil"
	"log"
	"net/http"
)

func GetWebSearch(keyword string) (string, error) {
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

func SearchFor(param string) {

}
