package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_process_CanRequest_false(t *testing.T) {
	// 避免 debugprint 输出
	temp := needDebug
	needDebug = false
	defer func() { needDebug = temp }()
	//
	ast := assert.New(t)
	p := new(process)
	p.requestTimestamp = newTimestamp(1, 1)
	ast.False(p.CanRequest())
}

func Test_process_CanRequest_true(t *testing.T) {
	// 避免 debugprint 输出
	temp := needDebug
	needDebug = false
	defer func() { needDebug = temp }()
	//
	ast := assert.New(t)
	p := new(process)
	p.requestTimestamp = nil
	ast.True(p.CanRequest())
}
