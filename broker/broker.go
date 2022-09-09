package main

import (
	// golang package
	"fmt"
	"log"
	"net/http"
	"time"
)

// main main.
func main() {
	subscriberList := map[string]map[string]bool{}
	queue := map[string][]string{}

	publishHandler := func(w http.ResponseWriter, req *http.Request) {
		topic := req.URL.Query().Get("topic")
		content := req.URL.Query().Get("message")

		fmt.Println("New Message Published in " + topic + " topic")

		if _, ok := subscriberList[topic]; ok {
			for key := range subscriberList[topic] {
				if _, ok := queue[key]; ok {
					queue[key] = append(queue[key], content)
				} else {
					queue[key] = []string{content}
				}
			}

		}
	}

	subscribeHandler := func(w http.ResponseWriter, req *http.Request) {
		host := req.URL.Query().Get("host")
		port := req.URL.Query().Get("port")
		topic := req.URL.Query().Get("topic")

		fmt.Println(host + ":" + port + " subsribed to " + topic)

		if _, ok := subscriberList[topic]; ok {
			subscriberList[topic][host+":"+port] = true
		} else {
			subscriberList[topic] = map[string]bool{host + ":" + port: true}
		}
	}

	unsubscribeHandler := func(w http.ResponseWriter, req *http.Request) {
		host := req.URL.Query().Get("host")
		port := req.URL.Query().Get("port")
		topic := req.URL.Query().Get("topic")

		fmt.Println(host + ":" + port + " unsubsribed from " + topic)

		if _, ok := subscriberList[topic]; ok {
			delete(subscriberList[topic], host+":"+port)
		}
	}

	ackHandler := func(w http.ResponseWriter, req *http.Request) {
		host := req.URL.Query().Get("host")
		port := req.URL.Query().Get("port")

		fmt.Println("Ack by " + host + ":" + port)

		if len(queue[host+":"+port]) > 0 {
			queue[host+":"+port] = queue[host+":"+port][1:]
		}
	}

	publishFunc := func() {
		for true {
			for sub, msgList := range queue {
				if len(msgList) > 0 {
					fmt.Println(queue[sub])
					fmt.Println("Sending to " + sub)
					http.Post("http://"+sub+"/receivemsg?msg="+msgList[0], "application/json", nil)
				}

			}
			time.Sleep(5 * time.Second)
		}
	}

	http.HandleFunc("/publish", publishHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/unsubscribe", unsubscribeHandler)
	http.HandleFunc("/ack", ackHandler)

	go publishFunc()

	log.Fatal(http.ListenAndServe(":8080", nil))

}
