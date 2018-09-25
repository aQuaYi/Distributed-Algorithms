package mutual

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_timestamp_String(t *testing.T) {
	ast := assert.New(t)
	ts := timestamp{
		time:    0,
		process: 1,
	}
	actual := ts.String()
	expected := "<T0:P1>"
	ast.Equal(expected, actual)

}
