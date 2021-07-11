package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func env_load() (string, string, string) {
	// .envファイルを読み込み環境変数とする
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}

	// 環境変数を読み込み
	localPath := os.Getenv("LOCAL_PATH")
	GCPKey := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	projectID := os.Getenv("PROJECT_ID")

	return localPath, GCPKey, projectID
}

func main() {
	// 環境変数の読み込み
	localPath, GCPKey, projectID := env_load()
	log.Print("Env-vars successfully loaded")
	log.Print("Backin' up files from", localPath, "to", projectID, "on GCP Storage…")

	// クライアントを建てる
	client := create_client(GCPKey)
	defer client.Close()
	log.Printf("successfully set a client in \"%s\"", projectID)

	// バケットを作成
	bucket, mes := create_bucket(*client, projectID)
	log.Print(mes)

	// バケットへファイルをコピー
	mes = copy_directory(*bucket, localPath)
	log.Print(mes)
}
