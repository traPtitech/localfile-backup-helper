package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// 環境変数を取得
	localPath, gcpKey, projectId, bucketName, storageClass, duration, bucketNum, webhookId, webhookSecret := EnvVarLoad()
	log.Print("Env-vars successfully loaded")

	log.Printf("Backin' up files from \"%s\" to \"%s\" on gcp Storage...", localPath, projectId)
	startTime := time.Now()

	// クライアントを建てる
	client, err := CreateClient(gcpKey, projectId)
	if err != nil {
		log.Print("Error: failed to load create client")
		panic(err)
	}
	defer client.Close()

	// bucketName + 月(mod n) をバケット名とし、bucketNameに再代入
	t := &startTime
	bucketName = fmt.Sprintf("%s-%d", bucketName, t.Month()%time.Month(bucketNum))

	// バケットを作成
	bucket, err := CreateBucket(*client, projectId, storageClass, duration, bucketName)
	if err != nil {
		if err.Error() == "Error 409" {
			log.Print("Bucket with the same name exists. Objects will be overwritten.")
		} else {
			log.Print("Error: failed to create bucket")
			panic(err)
		}
	}

	// バケットへディレクトリをコピー
	objectNum, err, errs := CopyDirectory(*bucket, localPath)
	if err != nil {
		log.Print("Error: failed to load local directory")
		panic(err)
	}
	log.Printf("%d file(s) successfully backed up, %d error(s) occured", objectNum, len(errs))
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
		log.Print("failed to send webhook")
		panic(err)
	}
}
