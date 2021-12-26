package main

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/golang/snappy"
	"golang.org/x/sync/semaphore"
	"google.golang.org/api/option"
)

const routineNum = 10 // 一度にファイルをコピーするgoroutineを実行できる数

var sem = semaphore.NewWeighted(routineNum) // goroutine実行数制御のためのセマフォ

func createClient(gcpKey string, projectId string) (*storage.Client, error) {
	// jsonで渡された鍵のサービスアカウントに紐づけられたクライアントを建てる
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcpKey))
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully set a client in \"%s\"", projectId)
	return client, err
}

func createBucket(client storage.Client, projectId string, storageClass string, duration int64, bucketName string) (*storage.BucketHandle, error) {
	// バケットとメタデータの設定
	bucket := client.Bucket(bucketName)
	bucketAttrs := &storage.BucketAttrs{
		StorageClass:      storageClass,
		Location:          "asia-northeast1",
		VersioningEnabled: true,
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
			log.Printf("Bucket \"%s\" exists. Objects will be overwritten.", bucketName)
			return bucket, nil
		} else {
			return nil, err
		}
	}

	log.Printf("Bucket \"%s\" successfully created", bucketName)
	return bucket, err
}

func copyDirectory(bucket storage.BucketHandle, localPath string) (int, error, []error) {
	var errs []error
	objectNum := 0

	// ローカルのディレクトリ構造を読み込み
	filePaths := []string{}
	err := filepath.Walk(localPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filePaths = append(filePaths, path)
		}
		return nil
	})
	if err != nil {
		return 0, err, nil
	}

	// 指定のディレクトリのファイルを並列処理で1つずつストレージにコピー(セマフォで一度に10個までに制限)
	for _, filePath := range filePaths {
		go func(filePath string) {
			if err := sem.Acquire(context.Background(), 1); err != nil {
				errs = append(errs, err)
				return
			}
			defer sem.Release(1)

			err = copyFile(bucket, filePath, strings.TrimPrefix(filePath, localPath+"/"))
			if err != nil {
				errs = append(errs, err)
			} else {
				objectNum++
			}
		}(filePath)
	}

	return objectNum, nil, errs
}

func copyFile(bucket storage.BucketHandle, filePath string, objectName string) error {
	// ローカルのファイルを開く
	original, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer original.Close()

	// 書き込むためのWriterを作成
	ctx := context.Background()
	writer := bucket.Object(objectName).NewWriter(ctx)
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
