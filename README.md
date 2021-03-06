# localfile-backup-helper

ローカルファイルのバックアップ用スクリプト  
ローカルストレージ上の任意のフラットなディレクトリからデータをコピーし、GCP Storage にバックアップします。  
※ 月一ペースで叩くことを想定したコードを想定しています。それ以外の場合は修正が必要です。  
※ 指定されたディレクトリ内にサブディレクトリがあった場合、挙動が予期しないものになる可能性があります。

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
  - `BUCKET_NAME`  
    バックアップ先のバケットの名前
    小文字・数字・記号が使えますが大文字が使えません
  - `STORAGE_CLASS`  
    データを格納するバケットのストレージクラス
  - `DURATION`  
    単一バージョンを保存する期間 (日数指定)
  - `PARALLEL_NUM`  
    同時にバックアップするファイル数
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
export BUCKET_NAME={name}
export STORAGE_CLASS={storage-class}
export DURATION={duration}
export PARALLEL_NUM={number}
export TRAQ_WEBHOOK_ID={id}
export TRAQ_WEBHOOK_SECRET={secret}

go run *.go
```
