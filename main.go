package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"modset/pkg/gcp"
	"modset/pkg/webhook"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func envLoad() (string, string, string, error) {
	// .envファイルを読み込み環境変数とする
	err := godotenv.Load()
	if err != nil {
		return "", "", "", err
	}

	// 環境変数を読み込み
	localPath := os.Getenv("LOCAL_PATH")
	gcpKey := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	projectId := os.Getenv("PROJECT_ID")

	log.Print("Env-vars successfully loaded")
	return localPath, gcpKey, projectId, err
}

func loadDir(localPath string) ([]fs.FileInfo, error) {
	// ローカルのディレクトリ構造を読み込み
	files, err := ioutil.ReadDir(localPath)
	if err != nil {
		return nil, err
	}

	return files, err
}

func main() {
	startTime := time.Now()

	// 環境変数の読み込み
	localPath, gcpKey, projectId, err := envLoad()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Backin' up files from", localPath, "to", projectId, "on gcp Storage…")

	// クライアントを建てる
	client, err := gcp.CreateClient(gcpKey, projectId)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// バケットを作成
	bucket, err := gcp.CreateBucket(*client, projectId)
	if err != nil {
		log.Fatal(err)
	}

	// ローカルのディレクトリ構造を読み込み
	files, err := loadDir(localPath)
	if err != nil {
		log.Fatal(err)
	}

	// バケットへファイルをコピー
	objectNum, errs := gcp.CopyDirectory(*bucket, files, localPath)
	log.Printf("%d file(s) successfully copied, %d error(s) occured", objectNum, len(errs))

	// Webhook用のメッセージを作成
	endTime := time.Now()
	buDuration := endTime.Sub(startTime)
	mes := webhook.CreateMes(startTime, buDuration, objectNum)
	log.Print(mes)
}
