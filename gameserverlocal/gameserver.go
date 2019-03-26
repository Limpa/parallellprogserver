package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	arraysize = 1000
)

var (
	highscores  []highscore
	activegames map[int]*game
)

type checkAnswerReq struct {
	Sum  int    `json:"sum"`
	Name string `json:"name"`
}

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

func init() {
	highscores = make([]highscore, 0)
	activegames = make(map[int]*game, 0)
	rand.Seed(time.Now().UTC().UnixNano())
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

func next(w http.ResponseWriter, r *http.Request) {
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
	time.Sleep(time.Millisecond * 5)
	w.Write([]byte(`{"number":` + strconv.Itoa(g.array[v]) + `}`))
}

func answer(w http.ResponseWriter, r *http.Request) {
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
	delete(activegames, t)
}

func main() {
	http.HandleFunc("/new", newGame)
	http.HandleFunc("/next", next)
	http.HandleFunc("/answer", answer)
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}

// helpers
func writeErrorResponse(w *http.ResponseWriter, err string, errcode int) {
	http.Error(*w, `{"error":"`+err+`"}`, errcode)
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
