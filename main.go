package main

import (
	"bytes"
	"net/http"
)

import "fmt"

func main() {
	var code string
	fmt.Print("Enter code: ")
	fmt.Scanf("%s", &code)
	var request bytes.Buffer
	request.WriteString("https://members.pumpingstationone.org/rfid/check/FrontDoor/")
	request.WriteString(code)
	resp, err := http.Get(request.String())
	if err != nil {
		fmt.Printf("Whoops!")
	}
	if resp.StatusCode == 200 {
		fmt.Println("Success!")
	} else if resp.StatusCode == 403 {
		fmt.Println("Membership status: Expired")
	} else {
		fmt.Println("Code not found")
	}

}
