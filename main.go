package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
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

func init() {
	// 環境変数を取得
	localPath = os.Getenv("LOCAL_PATH")
	gcpKey = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	projectId = os.Getenv("PROJECT_ID")
	bucketName = os.Getenv("BUCKET_NAME")
	storageClass = os.Getenv("STORAGE_CLASS")
	duration, _ = strconv.ParseInt(os.Getenv("DURATION"), 0, 64)
	webhookId = os.Getenv("TRAQ_WEBHOOK_ID")
	webhookSecret = os.Getenv("TRAQ_WEBHOOK_SECRET")

	// 環境変数がどれか一つでも空だったらエラーを吐いて終了
	if localPath == "" || gcpKey == "" || projectId == "" || bucketName == "" || storageClass == "" || duration == 0 || webhookId == "" || webhookSecret == "" {
		log.Print("Error: Failed to load env-vars")
		panic("empty env-var(s) exist")
	}
}

func main() {
	log.Print("Backin' up files from", localPath, "to", projectId, "on gcp Storage...")
	startTime := time.Now()

	// クライアントを建てる
	client, err := CreateClient()
	if err != nil {
		log.Print("Error: Failed to load create client")
		panic(err)
	}
	defer client.Close()

	// bucketName + バックアップ日 をバケット名とする
	t := &startTime
	bucketName := fmt.Sprintf("%s-%d-%d-%d", bucketName, t.Year(), t.Month(), t.Day())

	// バケットを作成
	bucket, err := CreateBucket(*client, bucketName)
	if err != nil {
		log.Print("Error: Failed to create bucket")
		panic(err)
	}

	// バケットへファイルをコピー
	objectNum, err, errs := CopyDirectory(*bucket)
	if err != nil {
		log.Print("Error: Failed to load local directory")
		panic(err)
	}
	log.Printf("%d file(s) successfully copied, %d error(s) occured", objectNum, len(errs))
	if len(errs) != 0 {
		for i, err := range errs {
			log.Printf("Error %d: %s", i, err)
		}
	}

	// Webhook用のメッセージを作成
	endTime := time.Now()
	buduration := endTime.Sub(startTime)
	mes := CreateMes(bucketName, startTime, buduration, objectNum, len(errs))

	// WebhookをtraQ Webhook Botに送信
	err = SendWebhook(mes)
	if err != nil {
		log.Print("Failed to send webhook")
		panic(err)
	}
}
