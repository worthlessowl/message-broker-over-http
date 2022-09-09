package main

import (
	// golang package
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// main main.
func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Publisher simulator")
	fmt.Println("---------------------")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')

		splitMsg := strings.Split(text[:len(text)-1], " ")

		if len(splitMsg) != 0 {
			if splitMsg[0] == "publish" {
				host := splitMsg[1]
				topic := splitMsg[2]
				msg := splitMsg[3]

				http.Post("http://"+host+"/publish?"+"topic="+topic+"&"+"message="+msg, "application/json", nil)
			}
		}

	}

}
