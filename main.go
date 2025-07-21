package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type JsonLogs struct {
	UserID       string `json:"user_id"`
	SessionID    string `json:"session_id"`
	Page         string `json:"page"`
	Event        string `json:"event"`
	Timestamp    int64  `json:"timestamp"`
	Device       string `json:"device"`
	Browser      string `json:"browser"`
	IP           string `json:"ip"`
	Location     string `json:"location"`
	ResponseTime int    `json:"response_time_ms"`
	Language     string `json:"language"`
}

var (
	Devices = []string{
		"mobile",
		"desktop",
		"tablet",
	}
	Browsers = []string{
		"chrome",
		"firefox",
		"safari",
		"edge",
	}
	Locations = []string{
		"US",
		"EU",
		"APAC",
	}
	Languages = []string{
		"en",
		"fr",
		"es",
		"de",
	}
)

var pages = []string{
	"/home",
	"/login",
	"/logout",
	"/product/123",
	"/product/456",
	"/cart",
	"/checkout",
	"/order-confirmation",
	"/profile",
	"/search?q=shoes",
	"/faq",
	"/support",
	"/terms",
	"/privacy-policy",
	"/contact-us",
	"/about-us",
	"/blog",
}

var events = []string{
	"page_view",
	"click_button",
	"add_to_cart",
	"remove_from_cart",
	"checkout_start",
	"checkout_complete",
	"login_success",
	"login_failure",
	"signup",
	"logout",
	"scroll",
	"hover",
	"search",
	"filter_applied",
	"form_submitted",
	"video_play",
	"video_pause",
}

func (j *JsonLogs) generateUserId() {
	j.UserID = "user_" + strconv.Itoa(rand.Intn(1000000))
}

func (j *JsonLogs) generateSessionId() {
	j.SessionID = "session_" + strconv.Itoa(rand.Intn(1000000))
}

func (j *JsonLogs) generatePage() {
	page := pages[rand.Intn(len(pages))]
	j.Page = page
}

func (j *JsonLogs) generateEvent() {
	event := events[rand.Intn(len(events))]
	j.Event = event
}

func (j *JsonLogs) generateIp() {
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
	j.IP = ip
}

func randomLog() JsonLogs {
	return JsonLogs{
		UserID:       fmt.Sprintf("user_%d", rand.Intn(1000000)),
		SessionID:    fmt.Sprintf("session_%d", rand.Intn(1000000)),
		Page:         pages[rand.Intn(len(pages))],
		Event:        events[rand.Intn(len(events))],
		Timestamp:    time.Now().Unix(),
		Device:       Devices[rand.Intn(len(Devices))],
		Browser:      Browsers[rand.Intn(len(Browsers))],
		IP:           fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256)),
		Location:     Locations[rand.Intn(len(Locations))],
		ResponseTime: rand.Intn(5000),
		Language:     Languages[rand.Intn(len(Languages))],
	}
}

func main() {
	//
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":          "0.0.0.0:9092",
		"linger.ms":                  5,       // Attend jusqu’à 5ms pour batcher les messages
		"batch.num.messages":         10000,   // Nombre max de messages dans un batch
		"queue.buffering.max.kbytes": 1048576, // 1 Go de buffer max
		"compression.type":           "lz4",   // Compression rapide et efficace
		"acks":                       "1",
		"statistics.interval.ms":     5000, // Envoie des stats toutes les 5 secondes
	})

	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} //else if ev.Headers != nil {
				// 	fmt.Printf("Message delivered: %s\n", ev.Value)
				// 	// fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
				// 	// 	*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				// } else if string(ev.Value) != "" && json.Valid(ev.Value) {
				// 	// Ça pourrait être les stats JSON
				// 	fmt.Printf("Stats JSON: %s\n", ev.Value)
				// }
			}
		}
	}()

	topic := "logs"
	numWorkers := 25 // nombre de goroutines producteurs
	messagesPerWorker := 100000

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < messagesPerWorker; i++ {
				j := randomLog()
				data, err := json.Marshal(j)
				if err != nil {
					log.Printf("[worker %d] Failed to marshal JSON: %s", workerID)
					continue
				}
				key := j.UserID + "_" + j.SessionID
				// Produce messages to Kafka topic
				msg := &kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
					Key:            []byte(key),
					Value:          data,
				}
				for {
					err := p.Produce(msg, nil)
					if err == nil {
						break
					}
					if err.(kafka.Error).Code() == kafka.ErrQueueFull {
						time.Sleep(1 * time.Millisecond)
					} else {
						log.Printf("[worker %d] Produce error: %v", workerID, err)
						break
					}
				}
			}
		}(w)
	}

	wg.Wait()

	p.Flush(15 * 1000) 
	p.Close()

}
