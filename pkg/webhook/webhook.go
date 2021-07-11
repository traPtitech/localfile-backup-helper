package webhook

// Webhook関連の処理を集めるパッケージ

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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
		`### traQローカルファイルのバックアップが保存されました
	バックアップ開始時刻: %s
	バックアップ所要時間: %f 分
	オブジェクト数: %d
	エラー数: %d`,
		startTime.Format(timeFormat), buDuration.Minutes(), objectNum, len(errs))

	log.Print("Webhook message generated")
	return mes
}

func SendWebhook(mes string) error {
	// リクエスト先url生成とメッセージの暗号化
	webhookUrl := "https://q.trap.jp/api/v3/webhooks/" + webhookId
	sig := calcHMACSHA1(mes)

	// リクエスト作成とヘッダーの設定
	req, err := http.NewRequest("POST", webhookUrl, strings.NewReader(mes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.Header.Set("X-TRAQ-Signature", sig)

	// クライアントからリクエストを送信、レスポンスを受ける
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// レスポンスの内容を確認
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	log.Printf("Sent webhook to traQ Webhook Bot (Status Code: %d, body: %s)", res.StatusCode, body)
	return err
}

func calcHMACSHA1(mes string) string {
	// メッセージをHMAC-SHA1でハッシュ化(Bot Consoleのコピペ)
	mac := hmac.New(sha1.New, []byte(webhookSecret))
	_, _ = mac.Write([]byte(mes))
	return hex.EncodeToString(mac.Sum(nil))
}
