package main

import (
	"db_technopark/application/server"
	"db_technopark/connection"
	"log"
)

func main() {
	conn, err := connection.InitDBConnection()
	if err != nil {
		log.Fatal("Can not connect to database: ", err)
	}
	s := server.NewServer(":8080", conn)
	log.Fatal(s.ListenAndServe())
}
