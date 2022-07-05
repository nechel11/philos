package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
	
	"golang.org/x/net/websocket"
)

const (
	num = 5
	minTime = 1000
	maxTime = 5000
)

func main() {
	handlefunc()
}

func handlefunc(){
	http.HandleFunc("/", index)
	http.Handle("/think", websocket.Handler(think))

	log.Print("Server runing on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "index.html")
	
}

func think(ws *websocket.Conn){
	output := make(chan string)
	forks := make([]sync.Mutex, num)
	for i:=1; i< num; i++{
		go eat(output, i, i-1, i, forks)
	}
	go eat(output, 0, num-1, 0, forks)
	for {
		websocket.Message.Send(ws, <-output)
	}
}

func eat (output chan string, l_f, r_f, id int, forks []sync.Mutex){
	for {
		t_eat := genTime()
		t_think:= genTime()

		output<- format("phil", id, "think")
		time.Sleep(time.Duration(t_think) * time.Millisecond)
		output<- format("phil", id, "wait")
		
		left_action := "right"
		right_action := "left"

		if id == 0{
			left_action = "left"
			right_action = "right"
		}

		forks[l_f].Lock()
		output<-format("fork", l_f, left_action) 
		forks[r_f].Lock()
		output<-format("fork", r_f, right_action) 

		output<-format("phil", id, "eat")
		time.Sleep(time.Duration(t_eat) * time.Millisecond)
		output<-format("phil", id, "wait")

		forks[r_f].Unlock()
		output<-format("fork", r_f, "back") 
		forks[l_f].Unlock()
		output<-format("fork", l_f, "back") 
	}
}

func genTime() int {
	return rand.Intn(maxTime+minTime) - minTime
}

func format(who string, id int, action string) string {
	return fmt.Sprintf(`{"who": "%s", "id": %d, "action": "%s"}`, who, id, action)
}
