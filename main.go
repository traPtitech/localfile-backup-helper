package main

import (
	"io/fs"
	"log"
	"os"
	"time"

	"github.com/traPtitech/localfile-backup-helper/gcp"
	"github.com/traPtitech/localfile-backup-helper/webhook"
)

var (
	localPath string
	projectId string
)

func init() {
	// 環境変数をグローバル変数に代入
	localPath = os.Getenv("LOCAL_PATH")
	projectId = os.Getenv("PROJECT_ID")

	if localPath == "" || projectId == "" {
		log.Print("Error: Failed to load env-vars")
		panic("empty env-var(s) exist")
	}
}

func loadDir() ([]fs.DirEntry, error) {
	// ローカルのディレクトリ構造を読み込み
	files, err := os.ReadDir(localPath)
	if err != nil {
		return nil, err
	}

	return files, err
}

func main() {
	log.Print("Backin' up files from", localPath, "to", projectId, "on gcp Storage...")
	startTime := time.Now()

	// クライアントを建てる
	client, err := gcp.CreateClient()
	if err != nil {
		log.Print("Error: Failed to load create client")
		panic(err)
	}
	defer client.Close()

	// バケットを作成
	bucket, err := gcp.CreateBucket(*client)
	if err != nil {
		log.Print("Error: Failed to create bucket")
		panic(err)
	}

	// ローカルのディレクトリ構造を読み込み
	files, err := loadDir()
	if err != nil {
		log.Print("Error: Failed to load local directory")
		panic(err)
	}

	// バケットへファイルをコピー
	objectNum, errs := gcp.CopyDirectory(*bucket, files)
	log.Printf("%d file(s) successfully copied, %d error(s) occured", objectNum, len(errs))
	if len(errs) != 0 {
		for i, err := range errs {
			log.Printf("Error %d: %s", i, err)
		}
	}

	// Webhook用のメッセージを作成
	endTime := time.Now()
	buDuration := endTime.Sub(startTime)
	mes := webhook.CreateMes(startTime, buDuration, objectNum, errs)

	// WebhookをtraQ Webhook Botに送信
	err = webhook.SendWebhook(mes)
	if err != nil {
		log.Print("Failed to send webhook")
		panic(err)
	}
}
