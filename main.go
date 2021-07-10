package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func env_load() (string, string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}
	localPath := os.Getenv("LOCAL_PATH")
	GCPKey := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	projectID := os.Getenv("PROJECT_ID")
	return localPath, GCPKey, projectID
}

func main() {
	localPath, GCPKey, projectID := env_load()
	log.Println("Backin' up files from", localPath, "to", projectID, "on GCP Storage …")
	client := create_client(GCPKey)
	defer client.Close()
	mes := create_bucket(distPath)
	log.Print(mes)
	mes = copy_file(localPath, distPath)
	log.Print(mes)
}
