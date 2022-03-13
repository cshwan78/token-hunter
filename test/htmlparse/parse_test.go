package test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/cshwan78/pkg/api/github"
)

func TestFoo(t *testing.T) {
	// todo test code
	data, err := ioutil.ReadFile("./testsample.html")
	if err != nil {
		t.Fatal(err)
	}
	htmlbody := string(data)

	result, err := github.ParseSearchHTML(htmlbody)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	for _, item := range result {
		fmt.Printf("item: %+v\n", item)
	}

}
