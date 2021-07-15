package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func init() {
	// 環境変数を取得
	EnvVarLoad()

	// 環境変数がどれか一つでも空だったらエラーを吐いて終了
	emptyVars := EnvVarEmptyCheck()
	if len(emptyVars) != 0 {
		log.Printf("Error: env-var(s) %s empty", strings.Join(emptyVars, ", "))
		panic("empty env-var(s) exist")
	}
}

func main() {
	// メモリー使用量のログを記録
	// defer profile.Start(profile.MemProfile, profile.ProfilePath("./prof"), profile.NoShutdownHook).Stop()

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
