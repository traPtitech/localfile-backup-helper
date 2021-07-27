package main

import (
	"log"
	"os"
	"strconv"
)

func EnvVarLoad() (string, string, string, string, string, int64, int, string, string) {
	// 環境変数を取得
	localPath := getEnv("LOCAL_PATH")
	gcpKey := getEnv("GOOGLE_APPLICATION_CREDENTIALS")
	projectId := getEnv("PROJECT_ID")
	bucketName := getEnv("BUCKET_NAME")
	storageClass := getEnv("STORAGE_CLASS")
	duration, _ := strconv.ParseInt(getEnv("DURATION"), 0, 64)
	bucketNum, _ := strconv.Atoi(getEnv("BUCKET_NUMBERS"))
	webhookId := getEnv("TRAQ_WEBHOOK_ID")
	webhookSecret := getEnv("TRAQ_WEBHOOK_SECRET")
	return localPath, gcpKey, projectId, bucketName, storageClass, duration, bucketNum, webhookId, webhookSecret
}

func getEnv(name string) string {
	// 指定された名前の環境変数を取得、空ならばエラーを吐いて終了
	loadedVar := os.Getenv(name)
	if loadedVar == "" {
		log.Printf("Error: env-var \"%s\" empty", name)
		panic("empty env-var exists")
	}

	return loadedVar
}
