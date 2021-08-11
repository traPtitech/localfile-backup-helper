package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func EnvVarLoad() ([]string, string, string, []string, string, int64, int, string, string) {
	// 環境変数を取得
	localPaths := getEnvSlice("LOCAL_PATH")
	gcpKey := getEnv("GOOGLE_APPLICATION_CREDENTIALS")
	projectId := getEnv("PROJECT_ID")
	bucketNames := getEnvSlice("BUCKET_NAME")
	storageClass := getEnv("STORAGE_CLASS")
	duration, _ := strconv.ParseInt(getEnv("DURATION"), 0, 64)
	bucketNum, _ := strconv.Atoi(getEnv("BUCKET_NUMBERS"))
	webhookId := getEnv("TRAQ_WEBHOOK_ID")
	webhookSecret := getEnv("TRAQ_WEBHOOK_SECRET")

	// localPathsとbuceketNamesの要素数が一致しない場合エラーを吐いて終了
	if len(localPaths) != len(bucketNames) {
		panic("Error: env_vars LOCAL_PATH and BUCKET_NAME doesn't have same number of elements")
	}

	log.Print("Env-vars successfully loaded")
	return localPaths, gcpKey, projectId, bucketNames, storageClass, duration, bucketNum, webhookId, webhookSecret
}

func getEnv(name string) string {
	// 指定された名前の環境変数を取得、空ならばエラーを吐いて終了
	loadedVar := os.Getenv(name)
	if loadedVar == "" {
		panic(fmt.Sprintf("Error: env-var \"%s\" empty", name))
	}

	return loadedVar
}

func getEnvSlice(name string) []string {
	// 指定された名前の環境変数を空白で区切りスライスで取得、空ならばエラーを吐いて終了
	loadedVar := strings.Split(os.Getenv(name), " ")
	if len(loadedVar) == 0 {
		panic(fmt.Sprintf("Error: env-var \"%s\" empty", name))
	}

	return loadedVar
}
