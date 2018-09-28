package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Message(t *testing.T) {
	ast := assert.New(t)
	//
	ts := newTimestamp(0, 0)
	m := newMessage(requestResource, 0, 0, OTHERS, ts)
	//
	expected := "{申请资源, Time:0, From:0, To:-1, <T0:P0>}"
	actual := m.String()
	ast.Equal(expected, actual)
	//
	m.msgType = releaseResource
	expected = "{释放资源, Time:0, From:0, To:-1, <T0:P0>}"
	actual = m.String()
	ast.Equal(expected, actual)
	//
	m.msgType = acknowledgment
	expected = "{确认收到, Time:0, From:0, To:-1, <T0:P0>}"
	actual = m.String()
	ast.Equal(expected, actual)
	//
}
