package communication

import "fmt"

func CheckNode() {
	for i := 0; i < 10; i++ {
		go sendRequest()
	}
}

func sendRequest() {
	fmt.Println("In a thread")
}
