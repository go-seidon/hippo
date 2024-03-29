package file

import (
	"fmt"

	"github.com/go-seidon/provider/datetime"
)

type UploadLocation interface {
	GetLocation() string
}

type dailyRotate struct {
	clock datetime.Clock
}

func (l *dailyRotate) GetLocation() string {
	currentTs := l.clock.Now()
	year := currentTs.Format("2006")
	month := currentTs.Format("01")
	day := currentTs.Format("02")
	return fmt.Sprintf("%s/%s/%s", year, month, day)
}

type DailyRotateParam struct {
	Clock datetime.Clock
}

func NewDailyRotate(p DailyRotateParam) *dailyRotate {
	var clock datetime.Clock
	if p.Clock != nil {
		clock = p.Clock
	} else {
		clock = datetime.NewClock()
	}

	l := &dailyRotate{
		clock: clock,
	}
	return l
}
