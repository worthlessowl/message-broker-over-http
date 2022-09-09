package main

import (
	// golang package
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// main main.
func main() {
	messages := []string{}
	target := "localhost:8080"
	host := "localhost"
	port := "7070"

	receiveHandler := func(w http.ResponseWriter, req *http.Request) {
		msg := req.URL.Query().Get("msg")

		messages = append(messages, msg)

		http.Post("http://"+target+"/ack?"+"host="+host+"&"+"port="+port, "application/json", nil)
	}

	http.HandleFunc("/receivemsg", receiveHandler)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Subscriber simulator")
	fmt.Println("---------------------")

	subsFunc := func() {
		for {
			fmt.Print("-> ")
			text, _ := reader.ReadString('\n')

			splitMsg := strings.Split(text[:len(text)-1], " ")

			if len(splitMsg) != 0 {
				if splitMsg[0] == "subscribe" {
					topic := splitMsg[1]

					http.Post("http://"+target+"/subscribe?"+"host="+host+"&"+"port="+port+"&"+"topic="+topic, "application/json", nil)

				}

				if splitMsg[0] == "messages" {
					for _, value := range messages {
						fmt.Println(value)
					}
				}
			}
		}
	}
	go subsFunc()

	log.Fatal(http.ListenAndServe(":7070", nil))
}
