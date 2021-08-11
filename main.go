package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// 環境変数を取得
	localPaths, gcpKey, projectId, bucketNames, storageClass, duration, bucketNum, webhookId, webhookSecret := EnvVarLoad()

	log.Printf("Backin' up files from %q to \"%s\" on gcp Storage...", localPaths, projectId)
	startTime := time.Now()

	// クライアントを建てる
	client, err := CreateClient(gcpKey, projectId)
	if err != nil {
		panic(fmt.Sprintf("Error: failed to load create client - %s", err))
	}
	defer client.Close()

	// localPathesの長さ(= bucketNamesの長さ)だけ処理を繰り返す
	for i := range localPaths {
		// bucketNamesのi番目 + 月(mod n) をバケット名とし、bucketNamesのi番目に再代入
		bucketNames[i] = fmt.Sprintf("%s-%d", bucketNames[i], startTime.Month()%time.Month(bucketNum))

		// バケットを作成
		bucket, err := CreateBucket(*client, projectId, storageClass, duration, bucketNames[i])
		if err != nil {
			panic(fmt.Sprintf("Error: failed to create bucket - %s", err))
		}

		// バケットへディレクトリをコピー
		objectNum, err, errs := CopyDirectory(*bucket, localPaths[i])
		if err != nil {
			panic(fmt.Sprintf("Error: failed to copy directory - %s", err))
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
		mes := CreateMes(localPaths[i], bucketNames[i], startTime, buduration, objectNum, len(errs))

		// WebhookをtraQ Webhook Botに送信
		err = SendWebhook(mes, webhookId, webhookSecret)
		if err != nil {
			panic(fmt.Sprintf("Error: failed to send webhook - %s", err))
		}
	}
}
