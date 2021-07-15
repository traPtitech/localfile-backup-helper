package main

import (
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
	localPath = os.Getenv("LOCAL_PATH")
	gcpKey = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	projectId = os.Getenv("PROJECT_ID")
	bucketName = os.Getenv("BUCKET_NAME")
	storageClass = os.Getenv("STORAGE_CLASS")
	duration, _ = strconv.ParseInt(os.Getenv("DURATION"), 0, 64)
	webhookId = os.Getenv("TRAQ_WEBHOOK_ID")
	webhookSecret = os.Getenv("TRAQ_WEBHOOK_SECRET")
}

func EnvVarEmptyCheck() []string {
	// 空の環境変数を全て入れたスライスを作成
	emptyVars := []string{}
	if localPath == "" {
		emptyVars = append(emptyVars, "LOCAL_PATH")
	}
	if gcpKey == "" {
		emptyVars = append(emptyVars, "GOOGLE_APPLICATION_CREDENTIALS")
	}
	if projectId == "" {
		emptyVars = append(emptyVars, "PROJECT_ID")
	}
	if bucketName == "" {
		emptyVars = append(emptyVars, "BUCKET_NAME")
	}
	if storageClass == "" {
		emptyVars = append(emptyVars, "STORAGE_CLASS")
	}
	if duration == 0 {
		emptyVars = append(emptyVars, "DURATION")
	}
	if webhookId == "" {
		emptyVars = append(emptyVars, "TRAQ_WEBHOOK_ID")
	}
	if webhookSecret == "" {
		emptyVars = append(emptyVars, "TRAQ_WEBHOOK_SECRET")
	}

	return emptyVars
}
