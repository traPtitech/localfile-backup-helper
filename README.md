# localfile-backup-helper

ローカルファイルのバックアップ用スクリプト  
ローカルストレージ上の任意のフラットなディレクトリからデータをコピーし、GCP Storage にバックアップします。  
指定されたディレクトリ内のサブディレクトリは検出されないので、バックアップはされません。

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
  - `STORAGECLASS`  
    データを格納するバケットのストレージクラス
  - `DURATION`
    データを格納する期間 (日数指定)
- traQ Webhook Bot 関連
  - `TRAQ_WEBHOOK_ID`  
    traQ Webhook Bot の ID
  - `TRAQ_WEBHOOK_SECRET`  
    traQ Webhook Bot のシークレット

## ローカルで動かす場合

シェルスクリプトで動かします。  
このリポジトリをクローンしたディレクトリ直下に下のような内容で任意の sh ファイルを作り、コンソールから`sh xxx.sh`で実行してください

```sh xxx.sh
#!/bin/sh

# 環境変数の設定
export LOCAL_PATH={path}
export GOOGLE_APPLICATION_CREDENTIALS={path}
export PROJECT_ID={project-id}
export TRAQ_WEBHOOK_ID={webhook-id}
export TRAQ_WEBHOOK_SECRET={secret}
export STORAGECLASS={strage-class}
export DURATION={duration}

go run main.go
```
