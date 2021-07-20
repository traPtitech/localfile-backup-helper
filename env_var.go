package main

import (
	"log"
	"os"
	"strconv"
)

// 環境変数をグローバル変数として定義
var (
	localPath     string
	gcpKey        string
	projectId     string
	bucketName    string
	storageClass  string
	duration      int64
	webhookId     string
	webhookSecret string
)

func EnvVarLoad() {
	// 環境変数を取得
	localPath = getEnv("LOCAL_PATH")
	gcpKey = getEnv("GOOGLE_APPLICATION_CREDENTIALS")
	projectId = getEnv("PROJECT_ID")
	bucketName = getEnv("BUCKET_NAME")
	storageClass = getEnv("STORAGE_CLASS")
	duration, _ = strconv.ParseInt(getEnv("DURATION"), 0, 64)
	webhookId = getEnv("TRAQ_WEBHOOK_ID")
	webhookSecret = getEnv("TRAQ_WEBHOOK_SECRET")
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
