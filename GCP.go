package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func create_client(GCPKey string) *storage.Client {
	ctx := context.Background()
	// jsonで渡された鍵のサービスアカウントに紐づけられたクライアントを建てる
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(GCPKey))
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func create_bucket(client storage.Client, projectID string) (*storage.BucketHandle, string) {
	// "s512_local" + バックアップ日時 をバケット名にする
	t := time.Now()
	bucketName := fmt.Sprintf("s512_local-%d-%d-%d", t.Year(), t.Month(), t.Day())

	// 10秒経ってもバケットが作成されない(タイムアウト)場合自動的にキャンセルするように設定
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// バケット名とメタデータの設定
	bucket := client.Bucket(bucketName)
	bucketAtters := &storage.BucketAttrs{
		StorageClass: "COLDLINE",
		Location:     "asia-northeast1",
		// 生成から90日でバケットを削除
		Lifecycle: storage.Lifecycle{Rules: []storage.LifecycleRule{
			{
				Action:    storage.LifecycleAction{Type: "Delete"},
				Condition: storage.LifecycleCondition{AgeInDays: 90},
			},
		}},
	}

	// バケットの作成
	err := bucket.Create(ctx, projectID, bucketAtters)
	if err != nil {
		log.Fatal(err)
	}

	return bucket, fmt.Sprintf("Bucket %s successfully created", bucketName)
}

func copy_file(localPath, distPath string) string {
	// ローカルのディレクトリ構造を読み込み
	bu_files, err := ioutil.ReadDir(localPath)
	if err != nil {
		log.Fatal(err)
	}

	// ファイルをストレージにコピー
	for _, file := range bu_files {
		copy, err := os.Create(distPath + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		original, err := os.Open(localPath + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(copy, original)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Copied", file.Name())
		copy.Close()
		original.Close()
	}
	return fmt.Sprintf("%d file(s) successfully copied", len(bu_files))
}
