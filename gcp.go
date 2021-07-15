package main

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/golang/snappy"
	"google.golang.org/api/option"
)

func CreateClient() (*storage.Client, error) {
	// jsonで渡された鍵のサービスアカウントに紐づけられたクライアントを建てる
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcpKey))
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully set a client in \"%s\"", projectId)
	return client, err
}

func CreateBucket(client storage.Client, bucketName string) (*storage.BucketHandle, error) {
	// バケットとメタデータの設定
	bucket := client.Bucket(bucketName)
	bucketAtters := &storage.BucketAttrs{
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
	err := bucket.Create(ctx, projectId, bucketAtters)
	if err != nil {
		return nil, err
	}

	log.Printf("Bucket \"%s\" successfully created", bucketName)
	return bucket, err
}

func CopyDirectory(bucket storage.BucketHandle) (int, error, []error) {
	var errs []error
	objectNum := 0

	// ローカルのディレクトリ構造を読み込み
	files, err := os.ReadDir(localPath)
	if err != nil {
		return 0, err, errs
	}

	// 指定のディレクトリのファイルを1つずつストレージにコピー
	for _, file := range files {
		err = copyFile(bucket, file)
		if err != nil {
			errs = append(errs, err)
		} else {
			objectNum++
		}
	}

	return objectNum, nil, errs
}

func copyFile(bucket storage.BucketHandle, file fs.DirEntry) error {
	// ローカルのファイルを開く
	original, err := os.Open(localPath + "/" + file.Name())
	if err != nil {
		return err
	}
	defer original.Close()

	// 書き込むためのWriterを作成
	ctx := context.Background()
	writer := bucket.Object(file.Name()).NewWriter(ctx)
	snappyWriter := snappy.NewBufferedWriter(writer)
	defer snappyWriter.Close()
	defer writer.Close()

	// 並列処理のためのパイプリーダー・ライターを定義
	pr, pw := io.Pipe()

	// 元のファイルをパイプライターに書き込み(パイプリーダーからGCP上ファイルへのコピーとの並行処理)
	errChan := make(chan error, 1)
	go func() {
		_, err := io.Copy(pw, original)
		defer pw.Close()
		errChan <- err
	}()

	// パイプリーダーを読んでGCP上のファイルに書きこみ
	_, err = io.Copy(snappyWriter, pr)
	if err != nil {
		return err
	}
	defer pr.Close()

	err = <-errChan
	if err != nil {
		return err
	}

	return err
}
