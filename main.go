package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func env_load() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}
	localPath := os.Getenv("LOCAL_PATH")
	distPath := os.Getenv("DIST_PATH")
	return localPath, distPath
}

func main() {
	localPath, distPath := env_load()
	log.Println("Backin' up files from", localPath, "to", distPath, "…")
	mes := create_bucket(distPath)
	log.Print(mes)
	mes = copy_file(localPath, distPath)
	log.Print(mes)
}
