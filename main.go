package main

import (
	"bytes"
	"net/http"
	"os"
)

import "fmt"

var code string

func main() {
	go http.HandleFunc("/", displayCode)
	go http.ListenAndServe(":8080", nil)
	for {
		fmt.Print("Enter code: ")
		fmt.Scanf("%s", &code)
		var request bytes.Buffer
		request.WriteString("https://members.pumpingstationone.org/rfid/check/FrontDoor/")
		request.WriteString(code)
		resp, err := http.Get(request.String())
		if err != nil {
			fmt.Printf("Whoops!")
			os.Exit(1)
		}
		if resp.StatusCode == 200 {
			fmt.Println("Success!")
			code = ""
		} else if resp.StatusCode == 403 {
			fmt.Println("Membership status: Expired")
		} else {
			fmt.Println("Code not found")
		}
	}

}

func displayCode(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(code))
}
