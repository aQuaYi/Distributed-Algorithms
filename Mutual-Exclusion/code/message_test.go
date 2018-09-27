package mutual

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Message(t *testing.T) {
	ast := assert.New(t)
	//
	ts := timestamp{time: 0, process: 0}
	m := newMessage2(requestResource, 0, others, ts)
	//
	expected := "{申请资源,<T0:P0>,From:0,To:-1}"
	actual := m.String()
	ast.Equal(expected, actual)
	//
	m.msgType = releaseResource
	expected = "{释放资源,<T0:P0>,From:0,To:-1}"
	actual = m.String()
	ast.Equal(expected, actual)
	//
	m.msgType = acknowledgment
	expected = "{确认收到,<T0:P0>,From:0,To:-1}"
	actual = m.String()
	ast.Equal(expected, actual)
	//
}
