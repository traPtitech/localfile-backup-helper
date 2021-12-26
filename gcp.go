package main

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/golang/snappy"
	"golang.org/x/sync/semaphore"
	"google.golang.org/api/option"
)

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

func copyDirectory(ctx context.Context, bucket storage.BucketHandle, localPath string, parallelNum int64) (int, error, []error) {
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

	sem := semaphore.NewWeighted(parallelNum) // goroutine実行数制御のためのセマフォ
	errs := []error{}
	errMutex := sync.Mutex{}
	objectNum := 0
	wg := sync.WaitGroup{}

	// 指定のディレクトリのファイルを並列処理で1つずつストレージにコピー(セマフォで一度に10個までに制限)
	for _, filePath := range filePaths {
		wg.Add(1)

		if err := sem.Acquire(ctx, 1); err != nil {
			errMutex.Lock()
			errs = append(errs, err)
			errMutex.Unlock()
		}

		go func(filePath string) {
			defer sem.Release(1)

			err = copyFile(bucket, filePath, strings.TrimPrefix(filePath, localPath+"/"))
			if err != nil {
				errs = append(errs, err)
			} else {
				objectNum++
			}

			wg.Done()
		}(filePath)
	}

	wg.Wait()
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
