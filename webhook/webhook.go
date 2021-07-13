// Webhook関連の処理を集めるパッケージ
package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strings"
)

// Webhookを扱う構造体の定義
type Handler struct {
	WebhookId     string
	WebhookSecret string
}

func (env *Handler) SendWebhook(mes string) error {
	// リクエスト先url生成とメッセージの暗号化
	webhookUrl := "https://q.trap.jp/api/v3/webhooks/" + env.WebhookId
	sig := env.calcHMACSHA1(mes)

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

func (env *Handler) calcHMACSHA1(mes string) string {
	// メッセージをHMAC-SHA1でハッシュ化(Bot Consoleのコピペ)
	mac := hmac.New(sha1.New, []byte(env.WebhookSecret))
	_, _ = mac.Write([]byte(mes))
	return hex.EncodeToString(mac.Sum(nil))
}
