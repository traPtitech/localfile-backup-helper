package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}
	localPath := os.Getenv("LOCAL_PATH")
	distPath := os.Getenv("DIST_PATH")
	fmt.Print(localPath, " ", distPath, "\n")
}
