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
	n_philos = 5
	min_time = 1000
	max_time = 5000
)

func main(){
	handle_func()
}

func handle_func(){
	http.HandleFunc("/", index)
	http.Handle("/think", websocket.Handler(think))

	log.Print("Server runing on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "index.html")
}

func think(ws *websocket.Conn){
	forks := make([]sync.Mutex, n_philos)
	chanel := make(chan string)

	for i := 1; i < n_philos; i++{
		go action(i, i-1, i, chanel, forks)
	}
	// first philo has to be left-handed, otherwise deadlock//
	go action(0, n_philos - 1, 0, chanel, forks) 
	for {
		websocket.Message.Send(ws, <-chanel)
	}
}

func action(left_fork, right_fork, id int, chanel chan string, forks []sync.Mutex){
	for {	
		
		t_eat := genTime()
		t_think := genTime()
		var leftaction string
		var rightaction string
		if_first(id, &leftaction, &rightaction)

		chanel <- format("phil", id, "think")
		time.Sleep(time.Duration(t_think) * time.Millisecond)
		chanel <- format("phil", id, "wait")

		/* philo takes forks and locks other routines */
		forks[left_fork].Lock()
		chanel <- format("fork" , left_fork, leftaction)
		forks[right_fork].Lock()
		chanel <- format("fork" , right_fork, rightaction)

		/* philo takes a dinner for t_eat milliseconds */
		chanel <- format("phil", id, "eat")
		time.Sleep(time.Duration(t_eat) * time.Millisecond)
		chanel <- format("phil", id, "wait")

		/* after having dinner philo unlocks forks (routines) */
		forks[right_fork].Unlock()
		chanel <- format("fork" , right_fork, "back")
		forks[left_fork].Unlock()
		chanel <- format("fork" , left_fork, "back")
	}
}

func if_first(id int, leftaction *string, rightaction *string){
	/* check if philos' id == 0 and set action */
	if id == 0{
		*leftaction = "left"
		*rightaction = "right"
	} else {
		*leftaction = "right"
		*rightaction = "left"
	}
}

func genTime() int {
	return rand.Intn(max_time+min_time) - min_time
}

func format(who string, id int, action string) string {
	return fmt.Sprintf(`{"who": "%s", "id": %d, "action": "%s"}`, who, id, action)
}
