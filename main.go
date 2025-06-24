package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type JsonLogs struct {
	UserID        string    `json:"user_id"`
	SessionID     string    `json:"session_id"`
	Page          string    `json:"page"`
	Event         string    `json:"event"`
	Timestamp     time.Time `json:"timestamp"`
	Device        string    `json:"device"`
	Browser       string    `json:"browser"`
	IP            string    `json:"ip"`
	Location      string    `json:"location"`
	ResponseTime  int       `json:"response_time_ms"`
	Language      string    `json:"language"`
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



func main() {
	// 
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})	
	
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
                } else {
                    fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
                        *ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
                }
            }
        }
	}()

	topic := "logs"
	
    for n := 0; n < 10; n++ {
		j := JsonLogs{}
		j.generateUserId()
		j.generateSessionId()
		j.generatePage()
		j.generateEvent()
		j.generateIp()
		j.Timestamp = time.Now()
		j.Device = Devices[rand.Intn(len(Devices))]
		j.Browser = Browsers[rand.Intn(len(Browsers))]
		j.Location = Locations[rand.Intn(len(Locations))]
		j.Language = Languages[rand.Intn(len(Languages))]
		j.ResponseTime = rand.Intn(5000) // Simulating response time in milliseconds
		data, err := json.Marshal(j)
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %s", err)
			os.Exit(1)
		}
		key := j.UserID + "_" + j.SessionID
		// Produce messages to Kafka topic
        p.Produce(&kafka.Message{
            TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
            Key:            []byte(key),
            Value:          []byte(data),
        }, nil)
    }

	p.Flush(15 * 1000) // Wait for messages to be delivered
	p.Close()


}
