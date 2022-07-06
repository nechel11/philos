package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
	"os"

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
		to_log("think", id, t_think)
		time.Sleep(time.Duration(t_think) * time.Millisecond)
		chanel <- format("phil", id, "wait")

		/* philo takes forks and locks other routines */
		forks[left_fork].Lock()
		to_log("left lock", id, 0)
		chanel <- format("fork" , left_fork, leftaction)
		forks[right_fork].Lock()
		to_log("right lock", id, 0)
		chanel <- format("fork" , right_fork, rightaction)

		/* philo takes a dinner for t_eat milliseconds */
		chanel <- format("phil", id, "eat")
		to_log("eat", id, t_eat)
		time.Sleep(time.Duration(t_eat) * time.Millisecond)
		chanel <- format("phil", id, "wait")

		/* after having dinner philo unlocks forks (routines) */
		forks[right_fork].Unlock()
		to_log("right unlock", id, 0)
		chanel <- format("fork" , right_fork, "back")
		forks[left_fork].Unlock()
		to_log("left unlock", id, 0)
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

func to_log(action string, id int, time int){
	f, err := os.OpenFile("logs", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	switch action {
		case "eat":
			log.Printf("Philo %d is eating for %d millisecs\n", id, time)
		case "think":
			log.Printf("Philo %d is thinking for %d millisecs\n", id, time)
		case "left lock":
			log.Printf("Philo %d locks left fork\n", id)
		case "right lock":
			log.Printf("Philo %d locks right fork\n", id)
		case "left unlock":
			log.Printf("Philo %d unlocks left fork\n", id)
		case "right unlock":
			log.Printf("Philo %d unlocks right fork\n", id)
	}
	
}

func format(who string, id int, action string) string {
	return fmt.Sprintf(`{"who": "%s", "id": %d, "action": "%s"}`, who, id, action)
}
