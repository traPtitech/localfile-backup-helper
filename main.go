package main

import (
	"io/fs"
	"log"
	"modset/pkg/gcp"
	"modset/pkg/webhook"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	localPath string
	projectId string
)

func envLoad() error {
	// .envファイルを読み込み環境変数とする
	err := godotenv.Load()
	if err != nil {
		return err
	}

	// mainで使う環境変数をグローバル変数に代入
	localPath = os.Getenv("LOCAL_PATH")
	projectId = os.Getenv("PROJECT_ID")

	// 付属パッケージに環境変数を設定
	gcp.EnvSet()
	webhook.EnvSet()

	log.Print("Env-vars successfully loaded")
	return err
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
	startTime := time.Now()

	// 環境変数の読み込み
	err := envLoad()
	if err != nil {
		log.Fatal("Failed to load env-vars: ", err)
	}

	log.Print("Backin' up files from", localPath, "to", projectId, "on gcp Storage…")

	// クライアントを建てる
	client, err := gcp.CreateClient()
	if err != nil {
		log.Fatal("Failed to load create client: ", err)
	}
	defer client.Close()

	// バケットを作成
	bucket, err := gcp.CreateBucket(*client)
	if err != nil {
		log.Fatal("Failed to create bucket: ", err)
	}

	// ローカルのディレクトリ構造を読み込み
	files, err := loadDir()
	if err != nil {
		log.Fatal("Failed to load local directory: ", err)
	}

	// バケットへファイルをコピー
	objectNum, errs := gcp.CopyDirectory(*bucket, files)
	log.Printf("%d file(s) successfully copied, %d error(s) occured", objectNum, len(errs))

	// Webhook用のメッセージを作成
	endTime := time.Now()
	buDuration := endTime.Sub(startTime)
	mes := webhook.CreateMes(startTime, buDuration, objectNum, errs)

	// WebhookをtraQ Webhook Botに送信
	err = webhook.SendWebhook(mes)
	if err != nil {
		log.Fatal("Failed to send webhook: ", err)
	}
}
