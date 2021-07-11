# traQ_local-backup-helper

s512 サーバーの traQ ローカルファイルのバックアップ用スクリプト  
ローカルの任意の場所からデータをコピーして GCP Storage にバックアップします。

## 設定

環境変数を必ず指定してください。

- ローカル関連
  - `LOCAL_PATH`  
    バックアップしたいローカルディレクトリのパス
- GCP 関連
  - `GOOGLE_APPLICATION_CREDENTIALS`  
    GCP の、バックアップ先のプロジェクトに紐づけられたサービスアカウントのキー(json ファイル)のパス
  - `PROJECT_ID`  
    バックアップ先のバケットを作成するプロジェクトの id
- traQ Webhook Bot 関連
  - `TRAQ_WEBHOOK_ID`  
    traQ Webhook Bot の ID
  - `TRAQ_WEBHOOK_SECRET`  
    traQ Webhook Bot のシークレット
