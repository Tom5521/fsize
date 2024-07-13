package stat_test

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/Tom5521/fsize/stat"
	"github.com/gookit/color"

	_ "unsafe"
)

//go:linkname parse github.com/Tom5521/fsize/stat.parseStatDate
func parse(string) (time.Time, error)

func TestParseDate(test *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			color.Errorln(r)
		}
	}()

	const (
		expectedInput  = "2024-06-16 21:01:08.044029927 -0400"
		expectedOutput = "2024-06-16 21:01:08"
	)

	t, err := parse(expectedInput)
	if err != nil {
		panic(err)
	}

	if t.Format(time.DateTime) != expectedOutput {
		fmt.Println("Time isn't parsed susseffully")
		test.Fail()
	}
	fmt.Println(t)
}
