package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func loadEnv() (string, string, string, string, string, int64, string, string) {
	// 環境変数を取得
	localPath := getEnv("LOCAL_PATH")
	gcpKey := getEnv("GOOGLE_APPLICATION_CREDENTIALS")
	projectID := getEnv("PROJECT_ID")
	bucketName := getEnv("BUCKET_NAME")
	storageClass := getEnv("STORAGE_CLASS")
	duration, _ := strconv.ParseInt(getEnv("DURATION"), 0, 64)
	webhookID := getEnv("TRAQ_WEBHOOK_ID")
	webhookSecret := getEnv("TRAQ_WEBHOOK_SECRET")

	log.Print("Env-vars successfully loaded")
	return localPath, gcpKey, projectID, bucketName, storageClass, duration, webhookID, webhookSecret
}

func getEnv(name string) string {
	// 指定された名前の環境変数を取得、空ならばエラーを吐いて終了
	loaded := os.Getenv(name)
	if loaded == "" {
		panic(fmt.Sprintf("Error: env-var \"%s\" empty", name))
	}

	return loaded
}
