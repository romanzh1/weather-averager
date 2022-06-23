package message

import (
	"fmt"
	"testing"
)

func BenchmarkGetWeatherAverage(b *testing.B) {
	fmt.Println(getWeatherAverage("дубна сегодня ave"))
}
