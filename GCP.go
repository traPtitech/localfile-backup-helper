package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func create_client(GCPKey string) *storage.Client {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(GCPKey))
	if err != nil {
		log.Fatal(err)
	}
	return client
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
		copy.Close()
		original.Close()
	}
	return fmt.Sprintf("%d file(s) successfully copied", len(bu_files))
}
