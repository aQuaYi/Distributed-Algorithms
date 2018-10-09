package mutualexclusion

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_debugPrintf_toPrint(t *testing.T) {
	rwm.Lock() // TODO: 删除此处的锁
	temp := needDebug
	needDebug = true
	rwm.Unlock()
	//
	var sb strings.Builder
	log.SetOutput(&sb)
	defer log.SetOutput(os.Stderr)
	//
	ast := assert.New(t)
	//
	words := "众鸟高飞尽，孤云独去闲。"
	//
	debugPrintf("%s", words)
	//
	ast.True(strings.Contains(sb.String(), words))
	// 还原 needDebug
	rwm.Lock() // TODO: 删除此处的锁
	needDebug = temp
	rwm.Unlock()
}

func Test_debugPrintf_notToPrint(t *testing.T) {
	rwm.Lock() // TODO: 删除此处的锁
	temp := needDebug
	needDebug = false
	rwm.Unlock()
	//
	var b strings.Builder
	log.SetOutput(&b)
	defer log.SetOutput(os.Stderr)
	//
	ast := assert.New(t)
	//
	words := "众鸟高飞尽，孤云独去闲。"
	//
	debugPrintf("%s", words)
	//
	ast.False(strings.Contains(b.String(), words))
	// 还原 needDebug
	rwm.Lock() // TODO: 删除此处的锁
	needDebug = temp
	rwm.Unlock()
}

func Test_max(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{

		{
			"a < b",
			args{
				1,
				2,
			},
			2,
		},

		{
			"a > b",
			args{
				2,
				1,
			},
			2,
		},

		{
			"a = b",
			args{
				2,
				2,
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := max(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("max() = %v, want %v", got, tt.want)
			}
		})
	}
}
