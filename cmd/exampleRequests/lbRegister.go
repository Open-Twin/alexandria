package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	lbUrl := "127.0.0.1:8080"
	data := url.Values{
		"ip": {"127.0.0.1"},
	}

	resp, err := http.PostForm("http://"+lbUrl+"/signup", data)
	if err != nil {
		fmt.Println("Error: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: %v", err)
	}

	if string(body) != "succesfully added" {
		fmt.Println("Adding node didn't work: %v", string(body))
	}
}
