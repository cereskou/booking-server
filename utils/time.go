package utils

import "time"

var (
	jst *time.Location
)

func init() {
	jst = time.FixedZone("Asia/Tokyo", 9*60*60)
}

//NowJST -
func NowJST() time.Time {
	return time.Now().UTC().In(jst)
}

//HourToSecond - 時間 to 秒
func HourToSecond(h int64) int64 {
	return h * 60 * 60
}
