package main

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/golang/snappy"
	"google.golang.org/api/option"
)

func CreateClient(gcpKey string, projectId string) (*storage.Client, error) {
	// jsonで渡された鍵のサービスアカウントに紐づけられたクライアントを建てる
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcpKey))
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully set a client in \"%s\"", projectId)
	return client, err
}

func CreateBucket(client storage.Client, projectId string, storageClass string, duration int64, bucketName string) (*storage.BucketHandle, error) {
	// バケットとメタデータの設定
	bucket := client.Bucket(bucketName)
	bucketAttrs := &storage.BucketAttrs{
		StorageClass: storageClass,
		Location:     "asia-northeast1",
		// 生成から90日でバケットを削除
		Lifecycle: storage.Lifecycle{Rules: []storage.LifecycleRule{
			{
				Action:    storage.LifecycleAction{Type: "Delete"},
				Condition: storage.LifecycleCondition{AgeInDays: duration},
			},
		}},
	}

	// バケットの作成
	ctx := context.Background()
	err := bucket.Create(ctx, projectId, bucketAttrs)
	if err != nil {
		// バケットが既にある場合のエラー(409: Conflict)を別枠で処理
		if strings.Contains(err.Error(), "Error 409") {
			log.Printf("Bucket \"%s\" already exists. Objects will be overwritten.", bucketName)
			return bucket, nil
		} else {
			return nil, err
		}
	}

	log.Printf("Bucket \"%s\" successfully created", bucketName)
	return bucket, err
}

func CopyDirectory(bucket storage.BucketHandle, localPath string) (int, error, []error) {
	var errs []error
	objectNum := 0

	// ローカルのディレクトリ構造を読み込み
	files, err := os.ReadDir(localPath)
	if err != nil {
		return 0, err, errs
	}

	// 指定のディレクトリのファイルを1つずつストレージにコピー
	for _, file := range files {
		err = copyFile(bucket, localPath+"/"+file.Name(), file.Name())
		if err != nil {
			errs = append(errs, err)
		} else {
			objectNum++
		}
	}

	return objectNum, nil, errs
}

func copyFile(bucket storage.BucketHandle, filePath string, fileName string) error {
	// ローカルのファイルを開く
	original, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer original.Close()

	// 書き込むためのWriterを作成
	ctx := context.Background()
	writer := bucket.Object(fileName).NewWriter(ctx)
	snappyWriter := snappy.NewBufferedWriter(writer)
	defer writer.Close()
	defer snappyWriter.Close()

	// GCP上のオブジェクトに書きこみ
	_, err = io.Copy(snappyWriter, original)
	if err != nil {
		return err
	}

	return err
}
