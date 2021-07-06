package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func env_setting() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}
	localPath := os.Getenv("LOCAL_PATH")
	distPath := os.Getenv("DIST_PATH")
	return localPath, distPath
}

func create_bucket(distPath string) string {
	t := time.Now()
	bucketName := fmt.Sprintf("s512_local-%d-%d-%d", t.Year(), t.Month(), t.Day())
	distPath += "/" + bucketName
	err := os.MkdirAll(distPath, 0777)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("Bucket %s successfully created", bucketName)
}

func copy_file(localPath string, distPath string) string {
	bu_files, err := ioutil.ReadDir(localPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range bu_files {
		copy, err := os.Create(distPath + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		original, err := os.Open(localPath + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(copy, original)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Copied", file.Name())
	}
	return fmt.Sprintf("%d file(s) successfully copied", len(bu_files))
}

func main() {
	localPath, distPath := env_setting()
	log.Println("Backin' up files from", localPath, "to", distPath, "…")
	mes := create_bucket(distPath)
	log.Print(mes)
	mes = copy_file(localPath, distPath)
	log.Print(mes)
}
