package message

import (
	"fmt"
	"testing"
)

func BenchmarkGetWeatherAverage(b *testing.B) {
	fmt.Println(GetWeatherAverage("дубна сегодня ave"))
}
