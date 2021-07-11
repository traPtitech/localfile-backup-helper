# server-backup

tokyotech.org、s512サーバーのtraQローカルファイルのバックアップスクリプト

ローカルからデータを抜いてGCP Storageにバックアップします。

## 設定
環境変数を必ず指定してください。

- ローカル
  - `LOCAL_PATH` バックアップしたいローカルディレクトリのパス
- GCP関連
  - `GOOGLE_APPLICATION_CREDENTIALS` GCPの、バックアップ先のプロジェクトに紐づけられたサービスアカウントのキー(jsonファイル)のパス
  - `PROJECT_ID` バックアップ先のバケットを作成するプロジェクトのid
- traQ関連
  - `TRAQ_WEBHOOK_ID` traQ Webhook BotのID
  - `TRAQ_WEBHOOK_SECRET` traQ Webhook Botのシークレット
