package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// 環境変数を取得
	localPath, gcpKey, projectId, bucketName, storageClass, duration, webhookId, webhookSecret := EnvVarLoad()

	log.Printf("Backin' up files from \"%s\" to \"%s\" on gcp Storage...", localPath, projectId)
	startTime := time.Now()

	// クライアントを建てる
	client, err := CreateClient(gcpKey, projectId)
	if err != nil {
		panic(fmt.Sprintf("Error: failed to load create client - %s", err))
	}
	defer client.Close()

	// バケットを作成
	bucket, err := CreateBucket(*client, projectId, storageClass, duration, bucketName)
	if err != nil {
		panic(fmt.Sprintf("Error: failed to create bucket - %s", err))
	}

	// バケットへディレクトリをコピー
	objectNum, err, errs := CopyDirectory(*bucket, localPath)
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
	buduration := endTime.Sub(startTime)
	mes := CreateMes(localPath, bucketName, startTime, buduration, objectNum, len(errs))

	// WebhookをtraQ Webhook Botに送信
	err = SendWebhook(mes, webhookId, webhookSecret)
	if err != nil {
		panic(fmt.Sprintf("Error: failed to send webhook - %s", err))
	}
}
