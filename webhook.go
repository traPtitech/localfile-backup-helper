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

func CreateMes(localPath string, bucketName string, startTime time.Time, buduration time.Duration, objectNum int, errs_num int) string {
	// traQに流すテキストメッセージを生成
	mes := fmt.Sprintf(
		`### ローカルファイルのバックアップが保存されました
	バックアップ元ディレクトリ: %s 
	生成/上書きされたバケット名: %s
	バックアップ開始時刻: %s
	バックアップ所要時間: %f 分
	オブジェクト数: %d
	エラー数: %d`,
		localPath, bucketName, startTime.Format(timeFormat), buduration.Minutes(), objectNum, errs_num)

	log.Print("Webhook message generated")
	return mes
}

func SendWebhook(mes string, webhookId string, webhookSecret string) error {
	// リクエスト先url生成とメッセージの暗号化
	webhookUrl := "https://q.trap.jp/api/v3/webhooks/" + webhookId
	sig := calcHMACSHA1(mes, webhookSecret)

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
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	log.Printf("Sent webhook to traQ Webhook Bot (Status Code: %d, body: %s)", res.StatusCode, body)
	return err
}

func calcHMACSHA1(mes string, webhookSecret string) string {
	// メッセージをHMAC-SHA1でハッシュ化(Bot Consoleのコピペ)
	mac := hmac.New(sha1.New, []byte(webhookSecret))
	_, _ = mac.Write([]byte(mes))
	return hex.EncodeToString(mac.Sum(nil))
}
