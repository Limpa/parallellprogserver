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

	req, err = http.NewRequest("GET", url+"/next", nil)
	if err != nil {
		log.Fatal(err)
	}
	token := b.Token
	req.Header.Set("X-Token", strconv.Itoa(token))
	sum := 0
	for i := 0; i < no_jobs; i++ {
		resp, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		b := struct {
			Number int `json:"number"`
		}{}
		json.Unmarshal(body, &b)
		sum += b.Number
	}

	req, err = http.NewRequest("POST", url+"/answer", bytes.NewBufferString(`{"name":"naive", "sum":`+strconv.Itoa(sum)+`}`))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Token", strconv.Itoa(token))
	resp, err = client.Do(req)
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Terminated with response: %s", string(body))
}
