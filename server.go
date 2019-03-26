package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"./gameserver"
	"./sandbox"
	"./wshub"
)

const (
	secret = "ff961dc5e8da688fa78540651160b223"
)

var (
	hub           *wshub.WSHub
	sbytes        []byte
	sendSolutions = false
	s             = [3]solution{
		{"naive", "https://gits-15.sys.kth.se/gist/linusri/47391dc55a3c5ad05052ce229b77637e"},
		{"better", "https://gits-15.sys.kth.se/gist/linusri/c619408d3f415356096e69882b2215f4"},
		{"worker", "https://gits-15.sys.kth.se/gist/linusri/7eeaca9a50b22a4cea3c9b16711b5955"},
	}
)

type solution struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

func solutions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if r.Header.Get("Auth") != secret {
			writeErrorResponse(&w, "provide secret in Auth header, hint MD5", 401)
			return
		}
		sendSolutions = true
		w.Write([]byte(`{"status":"success"}`))
	case "GET":
		if !sendSolutions {
			writeErrorResponse(&w, "forbidden", 403)
			return
		}
		w.Write(sbytes)
	default:
		writeErrorResponse(&w, "method not allowed", 405)
	}
}

func recieveFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeErrorResponse(&w, "method not allowed", 405)
		return
	}
	f, h, err := r.FormFile("file")
	if err != nil || f == nil || !strings.HasSuffix(h.Filename, ".go") {
		writeErrorResponse(&w, "invalid file", 400)
		return
	}
	sandbox.CreateJob(&f, h.Filename)
}

func onQueueChange() {
	queue := sandbox.GetQueue()
	hub.Broadcast(queue, "queue")
}

func init() {
	var err error
	sbytes, err = json.Marshal(struct {
		Solutions [3]solution `json:"solutions"`
	}{s})
	if err != nil {
		log.Fatal(err)
	}

	hub = wshub.New(gameserver.GetHighscores)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", hub.ConnectionHandler)
	http.HandleFunc("/solutions", solutions)
	http.HandleFunc("/upload", recieveFileUpload)
	go sandbox.Run(onQueueChange)
	go gameserver.Run(hub)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// helper
func writeErrorResponse(w *http.ResponseWriter, err string, errcode int) {
	http.Error(*w, `{"error":"`+err+`"}`, errcode)
}
