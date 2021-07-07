// GCP関係の処理をまとめたモジュール

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

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
