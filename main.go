package main

import (
	api "archive_lib/setup"
	"log"
	"os"
)

func main() {
	log.Println("PID:", os.Getpid())	
	api.Run()
}
