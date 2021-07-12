// GCPのライブラリを使った処理を集めるパッケージ
package gcp

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/golang/snappy"
	"google.golang.org/api/option"
)

var (
	localPath    string
	gcpKey       string
	projectId    string
	storageClass string
	duration     int64
)

func init() {
	// 環境変数をグローバル変数に代入
	localPath = os.Getenv("LOCAL_PATH")
	gcpKey = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	projectId = os.Getenv("PROJECT_ID")
	storageClass = os.Getenv("STORAGECLASS")
	duration, _ = strconv.ParseInt(os.Getenv("DURATION"), 0, 64)

	if localPath == "" || gcpKey == "" || projectId == "" || storageClass == "" || duration == 0 {
		log.Print("Error: Failed to load env-vars")
		panic("empty env-var(s) exist")
	}
}

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

func CreateBucket(client storage.Client) (*storage.BucketHandle, error) {
	// "localfile" + バックアップ日時 をバケット名にする
	t := time.Now()
	bucketName := fmt.Sprintf("localfile-%d-%d-%d", t.Year(), t.Month(), t.Day())

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
