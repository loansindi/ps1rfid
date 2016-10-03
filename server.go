package main

import (
	"net/http"
)

func serve(bot Robotter) {
	http.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Okay"))
		bot.openDoor()
	})
	http.ListenAndServe(":8080", nil)
}
