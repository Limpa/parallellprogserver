package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const arraysize = 1000

var (
	highscores  []highscore
	activegames map[int]*game
	sockets     []*websocket.Conn
	upgrader    websocket.Upgrader
)

type highscore struct {
	elapsedTime time.Duration
	name        string
}

type game struct {
	array     [arraysize]int
	startTime time.Time
	currval   int
	sum       int
	lock      sync.Mutex
}

func newGame(w http.ResponseWriter, r *http.Request) {
	id := rand.Int()
	var array [arraysize]int
	sum := 0
	for i := range array {
		array[i] = rand.Int() % 100
		sum += array[i]
	}
	activegames[id] = &game{array, time.Now(), 0, sum, sync.Mutex{}}
	w.Write([]byte(`{"token": ` + strconv.Itoa(id) + `}`))
}

func nextNumber(w http.ResponseWriter, r *http.Request) {
	g, _, err := getToken(r, &w)
	if err != nil {
		return
	}
	g.lock.Lock()
	v := g.currval
	g.currval = g.currval + 1
	g.lock.Unlock()
	if v >= len(g.array) {
		writeErrorResponse(&w, "no numbers left", 400)
		return
	}
	w.Write([]byte(`{"number":` + strconv.Itoa(g.array[v]) + `}`))
}

func writeErrorResponse(w *http.ResponseWriter, err string, errcode int) {
	http.Error(*w, `{"error":"`+err+`"}`, errcode)
}

type checkAnswerReq struct {
	Sum  int    `json:"sum"`
	Name string `json:"name"`
}

func getToken(r *http.Request, w *http.ResponseWriter) (*game, int, error) {
	token, err := strconv.Atoi(r.Header.Get("X-Token"))
	if err != nil {
		writeErrorResponse(w, "invalid token", 400)
		return nil, 0, err
	}
	g, ok := activegames[token]
	if !ok {
		writeErrorResponse(w, "invalid token", 400)
		return nil, 0, errors.New("error")
	}
	return g, token, nil
}

func checkAnswer(w http.ResponseWriter, r *http.Request) {
	g, t, err := getToken(r, &w)
	if err != nil {
		return
	}
	defer r.Body.Close()
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(&w, "error reading request body", 500)
		return
	}
	req := checkAnswerReq{}
	json.Unmarshal(payload, &req)
	if req.Sum != g.sum {
		writeErrorResponse(&w, "the sum isnt correct", 400)
		return
	}
	elapsedt := time.Since(g.startTime)
	w.Write([]byte(`{"message": "success", "time": "` + elapsedt.String() + `"}`))
	highscores = append(highscores, highscore{elapsedt, req.Name})
	hs := getHighscoreByteArray()
	for _, s := range sockets {
		s.WriteMessage(websocket.TextMessage, hs)
	}
	delete(activegames, t)
}

func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, getHighscoreByteArray())
	if err != nil {
		fmt.Println(err)
	}
	sockets = append(sockets, conn)
}

func getHighscoreByteArray() []byte {
	sort.Slice(highscores, func(i, j int) bool {
		return highscores[i].elapsedTime.Nanoseconds() < highscores[j].elapsedTime.Nanoseconds()
	})
	var buffer bytes.Buffer
	buffer.WriteString(`{"highscores":[`)
	for i, hs := range highscores {
		buffer.WriteString(fmt.Sprintf(`{"name":"%s", "time":"%s"}`, hs.name, hs.elapsedTime.String()))
		if i < len(highscores)-1 {
			buffer.WriteString(`,`)
		}
	}
	buffer.WriteString("]}")
	return buffer.Bytes()
}

func main() {
	highscores = make([]highscore, 0)
	activegames = make(map[int]*game, 0)
	rand.Seed(time.Now().UTC().UnixNano())
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/new", newGame)
	http.HandleFunc("/next", nextNumber)
	http.HandleFunc("/answer", checkAnswer)
	http.HandleFunc("/ws", ws)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
