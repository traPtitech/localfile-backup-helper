// GCPのライブラリを使った処理を集めるパッケージ
package gcp

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

// 環境変数を管理する構造体の定義
type GcpEnv struct {
	LocalPath    string
	GcpKey       string
	ProjectId    string
	StorageClass string
	Duration     int64
}

func (env *GcpEnv) CreateClient() (*storage.Client, error) {
	// jsonで渡された鍵のサービスアカウントに紐づけられたクライアントを建てる
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(env.GcpKey))
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully set a client in \"%s\"", env.ProjectId)
	return client, err
}

func (env *GcpEnv) CreateBucket(client storage.Client, bucketName string) (*storage.BucketHandle, error) {
	// バケットとメタデータの設定
	bucket := client.Bucket(bucketName)
	bucketAtters := &storage.BucketAttrs{
		StorageClass: env.StorageClass,
		Location:     "asia-northeast1",
		// 生成から90日でバケットを削除
		Lifecycle: storage.Lifecycle{Rules: []storage.LifecycleRule{
			{
				Action:    storage.LifecycleAction{Type: "Delete"},
				Condition: storage.LifecycleCondition{AgeInDays: env.Duration},
			},
		}},
	}

	// バケットの作成
	ctx := context.Background()
	err := bucket.Create(ctx, env.ProjectId, bucketAtters)
	if err != nil {
		return nil, err
	}

	log.Printf("Bucket \"%s\" successfully created", bucketName)
	return bucket, err
}

func (env *GcpEnv) CopyDirectory(bucket storage.BucketHandle) (int, error, []error) {
	var errs []error
	objectNum := 0

	// ローカルのディレクトリ構造を読み込み
	files, err := os.ReadDir(env.LocalPath)
	if err != nil {
		return 0, err, errs
	}

	// 指定のディレクトリのファイルを1つずつストレージにコピー
	for _, file := range files {
		err = env.copyFile(bucket, file)
		if err != nil {
			errs = append(errs, err)
		} else {
			objectNum++
		}
	}

	return objectNum, nil, errs
}

func (env *GcpEnv) copyFile(bucket storage.BucketHandle, file fs.DirEntry) error {
	// ローカルのファイルを開く
	original, err := os.Open(env.LocalPath + "/" + file.Name())
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
