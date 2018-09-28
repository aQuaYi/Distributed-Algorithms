package main

import (
	"testing"

	"github.com/aQuaYi/observer"
	"github.com/stretchr/testify/assert"
)

func Test_process_needResource_true(t *testing.T) {
	// 避免 debugprint 输出
	temp := needDebug
	needDebug = false
	defer func() { needDebug = temp }()
	//
	ast := assert.New(t)
	p := newProcess(10, 1, nil, observer.NewProperty(1))
	p.addOccupyTimes(1)
	ast.True(p.needResource())
}

func Test_process_needResource_false(t *testing.T) {
	// 避免 debugprint 输出
	temp := needDebug
	needDebug = false
	defer func() { needDebug = temp }()
	//
	ast := assert.New(t)
	p := newProcess(10, 1, nil, observer.NewProperty(1))
	ast.False(p.needResource())
}
