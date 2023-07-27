package main

import (
	"fmt"
	"log"
	"net/url"
	"testing"
)

func TestFileServer(t *testing.T) {
	// fmt.Println("_______")
	// files := http.FileServer(http.Dir("/public"))
	// fmt.Println(files)
	// http.Handle("/", files)
	// http.ListenAndServe(":8080", nil)
	// fmt.Println("_______")

	u, err := url.Parse("http://bing.com/search")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u.Query())
	fmt.Println(u.Query().Get("id"))
	fmt.Printf("%T", u.Query().Get("id"))
}
