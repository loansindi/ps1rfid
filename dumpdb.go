package main

import(
	"log"
	"fmt"
	"github.com/boltdb/bolt"
)

func main() {
	var err error

	fmt.Println("Opening database...")
	cacheDB, err := bolt.Open("rfid-tags.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cacheDB.Close()

	fmt.Println("Database opened, now reading entries")
	tx, _ := cacheDB.Begin(false)
	cursor := tx.Bucket([]byte("RFIDBucket")).Cursor()

	for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
		fmt.Printf("tag = %s, value = %s\n", k, v)
	}
	
	fmt.Println("Done!")
}
