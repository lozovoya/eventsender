package main

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type CallEventDTO struct {
	EventID   string `json:"event_id"`
	Type      string `json:"event_type"`
	Timestamp string `json:"timestamp"`
	Data      Call   `json:"data"`
}

type Call struct {
	CallID        string `json:"call_id"`
	Queue_ID      string `json:"queue_id"`
	DialedNumber  string `json:"dialed_phone"`
	CallingNumber string `json:"calling_number"`
	CallingLevel  string `json:"calling_level"`
}

func main() {
	queues := map[int]string{
		0: "queue0",
		1: "queue1",
		2: "queue2",
		3: "queue3",
		4: "queue4",
		5: "queue5",
		6: "queue6",
		7: "queue7",
		8: "queue8",
		9: "queue9",
	}

	levels := map[int]string{
		0: "Gold",
		1: "Silver",
		2: "Bronze",
	}
	for {

		n := 300
		for i := 1; i < n; i++ {
			//wg := sync.WaitGroup{}
			go func() {

				rand.Seed(time.Now().UnixNano())
				pause := rand.Intn(10)
				time.Sleep(time.Second * time.Duration(pause))
				var event CallEventDTO
				event.EventID = uuid.New().String()
				event.Type = "OnQueueInEvent"
				event.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
				event.Data.Queue_ID = queues[rand.Intn(len(queues))]
				event.Data.CallID = uuid.New().String()
				event.Data.DialedNumber = strconv.FormatInt(rand.Int63n(1000000)+1000, 10)
				event.Data.CallingLevel = levels[rand.Intn(len(levels))]
				event.Data.CallingNumber = strconv.FormatInt(rand.Int63n(1000000)+1000, 10)

				data, err := json.Marshal(event)
				if err != nil {
					log.Println(err)
					return
				}
				client := http.Client{Timeout: time.Second * 5}
				_, err = client.Post("http://localhost:9999/api/v1/events",
					"application/json",
					bytes.NewBuffer(data))

				pause = rand.Intn(60)
				time.Sleep(time.Second * time.Duration(pause))
				event.Type = "OnQueueOutEvent"
				event.EventID = uuid.New().String()
				data, err = json.Marshal(event)
				if err != nil {
					log.Println(err)
					return
				}
				_, err = client.Post("http://localhost:9999/api/v1/events",
					"application/json",
					bytes.NewBuffer(data))
				if err != nil {
					log.Println(err)
				}
			}()
		}

		log.Println("one turn finished")
		time.Sleep(time.Second * 15)

	}
}
