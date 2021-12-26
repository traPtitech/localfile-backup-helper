package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const timeFormat = "2006/01/02 15:04:05"

func createMes(localPath string, bucketName string, startTime time.Time, buDuration time.Duration, objectNum int, errs_num int) string {
	// traQに流すテキストメッセージを生成
	mes := fmt.Sprintf(
		`### ローカルファイルのバックアップが保存されました
	バックアップ元ディレクトリ: %s
	生成/上書きされたバケット名: %s
	バックアップ開始時刻: %s
	バックアップ所要時間: %f 分
	オブジェクト数: %d
	エラー数: %d`,
		localPath, bucketName, startTime.Format(timeFormat), buDuration.Minutes(), objectNum, errs_num)

	log.Print("Webhook message generated")
	return mes
}

func sendWebhook(mes string, webhookId string, webhookSecret string) error {
	if webhookId == "" || webhookSecret == "" {
		log.Printf("As webhook ID and secret being empty, webhook message will not be sent.\nHere is the message to be sent:\n%s", mes)
		return nil
	}

	// リクエスト先url生成とメッセージの暗号化
	webhookUrl := "https://q.trap.jp/api/v3/webhooks/" + webhookId
	sig := calcHash(mes, webhookSecret)

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

	// レスポンスの内容を確認
	if res.StatusCode >= 300 {
		return fmt.Errorf("status code: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	log.Printf("Sent webhook to traQ Webhook Bot (Status Code: %d, body: %s)", res.StatusCode, body)
	return err
}

func calcHash(mes string, webhookSecret string) string {
	// メッセージをHMAC-SHA1でハッシュ化(Bot Consoleのコピペ)
	mac := hmac.New(sha1.New, []byte(webhookSecret))
	_, _ = mac.Write([]byte(mes))
	return hex.EncodeToString(mac.Sum(nil))
}
