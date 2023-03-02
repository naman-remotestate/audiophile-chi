package main

import (
	"audiophile/database"
	"audiophile/server"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {

	done := make(chan int)
	srv := server.SetUpRoutes()

	if err := database.ConnectToDatabase(
		"localhost",
		"5433",
		"audiophile",
		"local",
		"local",
		database.SSLModeDisable); err != nil {
		logrus.Panicf("Failed to initialize and migrate database with error: %+v", err)
	}
	logrus.Print("migration successful!!")
	go func() {
		err := srv.Start(":3000")
		if err != nil {
			log.Fatalf("Failed to run server with error: %+v", err)
		}

	}()

	fmt.Println("Server Running at port :3000")
	<-done
}
