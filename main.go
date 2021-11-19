package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// 環境変数を取得
	localPath, gcpKey, projectID, bucketName, storageClass, duration, webhookID, webhookSecret := loadEnv()

	log.Printf("Backin' up files from \"%s\" to \"%s\" on gcp Storage...", localPath, projectID)
	startTime := time.Now()

	// クライアントを建てる
	client, err := createClient(gcpKey, projectID)
	if err != nil {
		panic(fmt.Sprintf("Error: failed to load create client - %s", err))
	}
	defer client.Close()

	// バケットを作成
	bucket, err := createBucket(*client, projectID, storageClass, duration, bucketName)
	if err != nil {
		panic(fmt.Sprintf("Error: failed to create bucket - %s", err))
	}

	// バケットへディレクトリをコピー
	objectNum, err, errs := copyDirectory(*bucket, localPath)
	if err != nil {
		panic(fmt.Sprintf("Error: failed to copy directory - %s", err))
	}
	log.Printf("%d file(s) successfully backed up, %d error(s) occurred", objectNum, len(errs))
	if len(errs) != 0 {
		for i, err := range errs {
			log.Printf("Error %d: %s", i, err)
		}
	}

	// Webhook用のメッセージを作成
	endTime := time.Now()
	buDuration := endTime.Sub(startTime)
	mes := createMes(localPath, bucketName, startTime, buDuration, objectNum, len(errs))

	// WebhookをtraQ Webhook Botに送信
	err = sendWebhook(mes, webhookID, webhookSecret)
	if err != nil {
		panic(fmt.Sprintf("Error: failed to send webhook - %s", err))
	}
}
