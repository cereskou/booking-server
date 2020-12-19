package utils

import "time"

var (
	jst *time.Location
)

//JST -
type JST struct {
	tm time.Time
}

func init() {
	jst = time.FixedZone("Asia/Tokyo", 9*60*60)
}

//NowJST -
func NowJST() *JST {
	return &JST{
		tm: time.Now().UTC().In(jst),
	}
}

//String -
func (t *JST) String() string {
	return t.tm.Format("2006/01/02 15:04:05")
}

//Time -
func (t *JST) Time() time.Time {
	return t.tm
}

//HourToSecond - 時間 to 秒
func HourToSecond(h int64) int64 {
	return h * 60 * 60
}
