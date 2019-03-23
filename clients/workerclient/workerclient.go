package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	url     = "http://localhost:8080"
	workers = 50
	no_jobs = 1000
)

func worker(jobs chan int, sum chan int, token string) {
	req, err := http.NewRequest("GET", url+"/next", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Token", token)

	client := http.Client{}
	rstruct := struct{ Number int }{}
	for _ = range jobs {
		resp, err := client.Do(req)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(body, &rstruct)
		sum <- rstruct.Number
	}
}

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url+"/new", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	b := struct {
		Token int `json:"token"`
	}{}
	json.Unmarshal(body, &b)
	token := strconv.Itoa(b.Token)

	results := make(chan int, no_jobs)
	jobs := make(chan int, no_jobs)
	for i := 0; i < workers; i++ {
		go worker(jobs, results, token)
	}
	for i := 0; i < no_jobs; i++ {
		jobs <- i
	}
	close(jobs)

	sum := 0
	for i := 0; i < no_jobs; i++ {
		sum += <-results
	}

	req, err = http.NewRequest("POST", url+"/answer", bytes.NewBufferString(`{"name":"worker", "sum":`+strconv.Itoa(sum)+`}`))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Token", token)
	resp, err = client.Do(req)
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("6")
	}

	fmt.Printf("Terminated with response: %s", string(body))
}
