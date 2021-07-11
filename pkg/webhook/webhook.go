package webhook

import (
	"fmt"
	"log"
	"os"
	"time"
)

const timeFormat = "2006/01/02 15:04:05"

var (
	webhookId     string
	webhookSecret string
)

func EnvSet() {
	// 環境変数をグローバル変数に代入
	webhookId = os.Getenv("TRAQ_WEBHOOK_ID")
	webhookSecret = os.Getenv("TRAQ_WEBHOOK_SECRET")
}

func CreateMes(startTime time.Time, buDuration time.Duration, objectNum int, errs []error) string {
	// traQに流すテキストメッセージを生成
	mes := fmt.Sprintf(
		`### s512ローカルファイルのバックアップが保存されました
	バックアップ開始時刻: %s
	バックアップ所要時間: %f 分
	オブジェクト数: %d
	エラー数: %d`,
		startTime.Format(timeFormat), buDuration.Minutes(), objectNum, len(errs))

	return mes
}

func SendWebhook(mes string) error {
	webhookUrl := "https://q.trap.jp/api/v3/webhooks/" + webhookId
	log.Print(webhookUrl, webhookSecret)
	return nil
}
