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

const url = "http://localhost:9090"
const no_jobs = 1000

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url+"/new", nil)
	if err != nil {
		log.Fatal("1")
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("2")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("3")
	}
	b := struct {
		Token int `json:"token"`
	}{}
	json.Unmarshal(body, &b)
	token := strconv.Itoa(b.Token)

	ch := make(chan int, no_jobs)
	for i := 0; i < no_jobs; i++ {
		go sendReq(ch, token)
	}
	sum := 0
	for i := 0; i < no_jobs; i++ {
		sum += <-ch
	}

	req, err = http.NewRequest("POST", url+"/answer", bytes.NewBufferString(`{"name":"better", "sum":`+strconv.Itoa(sum)+`}`))
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

func sendReq(ch chan int, token string) {
	req, err := http.NewRequest("GET", url+"/next", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Token", token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("8")
	}
	b := struct {
		Number int `json:"number"`
	}{}
	json.Unmarshal(body, &b)
	ch <- b.Number
}
