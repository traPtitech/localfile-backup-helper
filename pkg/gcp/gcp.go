package gcp

// GCPのライブラリを使った処理を集めるパッケージ

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/golang/snappy"
	"google.golang.org/api/option"
)

var (
	localPath string
	gcpKey    string
	projectId string
)

func EnvSet() {
	// 環境変数をグローバル変数に代入
	localPath = os.Getenv("LOCAL_PATH")
	gcpKey = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	projectId = os.Getenv("PROJECT_ID")
}

func CreateClient() (*storage.Client, error) {
	ctx := context.Background()

	// jsonで渡された鍵のサービスアカウントに紐づけられたクライアントを建てる
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcpKey))
	if err != nil {
		return nil, err
	}

	log.Printf("successfully set a client in \"%s\"", projectId)
	return client, err
}

func CreateBucket(client storage.Client) (*storage.BucketHandle, error) {
	// "s512_local" + バックアップ日時 をバケット名にする
	t := time.Now()
	bucketName := fmt.Sprintf("s512_local-%d-%d-%d", t.Year(), t.Month(), t.Day())

	// バケットとメタデータの設定
	ctx := context.Background()
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
	err := bucket.Create(ctx, projectId, bucketAtters)
	if err != nil {
		return nil, err
	}

	log.Printf("Bucket \"%s\" successfully created", bucketName)
	return bucket, err
}

func CopyDirectory(bucket storage.BucketHandle, files []fs.FileInfo) (int, []error) {
	var errs []error
	objectNum := 0

	// 指定のディレクトリのファイルを1つずつストレージにコピー
	for _, file := range files {
		err := copyFile(bucket, file)
		if err != nil {
			errs = append(errs, err)
		} else {
			objectNum++
		}
	}

	return objectNum, errs
}

func copyFile(bucket storage.BucketHandle, file fs.FileInfo) error {
	// ローカルのファイルを開く
	original, err := os.Open(localPath + "/" + file.Name())
	if err != nil {
		return err
	}
	defer original.Close()

	//書き込むためのWriterを作成
	ctx := context.Background()
	writer := bucket.Object(file.Name()).NewWriter(ctx)
	snappyWriter := snappy.NewBufferedWriter(writer)
	defer snappyWriter.Close()
	defer writer.Close()

	// 書きこみ
	_, err = io.Copy(snappyWriter, original)
	if err != nil {
		return err
	}

	return err
}
