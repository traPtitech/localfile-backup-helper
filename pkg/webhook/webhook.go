package webhook

import (
	"fmt"
	"time"
)

const dateFormat = "2006/01/02 15:04:05"

func CreateMes(startTime time.Time, buDuration time.Duration, objectNum int, errs []error) string {
	mes := fmt.Sprintf(`### s512ローカルファイルのバックアップが保存されました
	バックアップ開始時刻: %s
	バックアップ所要時間: %f 分
	オブジェクト数: %d
	エラー数: %d`,
		startTime.Format(dateFormat), buDuration.Minutes(), objectNum, len(errs))
	return mes
}
