package mutual

import (
	"testing"
)

func Test_start(t *testing.T) {
	type args struct {
		size         int
		occupyNumber int
		resource     *resource
	}
	tests := []struct {
		name string
		args args
	}{

		{
			"3 process, 9 occupy",
			args{
				size:         3,
				occupyNumber: 9,
				resource:     newResource(),
			},
		},

		// {
		// 	"9 process, 99 occupy",
		// 	args{
		// 		size:         9,
		// 		occupyNumber: 99,
		// 		resource:     newResource(),
		// 	},
		// },

		// {
		// 	"6 process, 999 occupy",
		// 	args{
		// 		size:         6,
		// 		occupyNumber: 999,
		// 		resource:     newResource(),
		// 	},
		// },

		// {
		// 	"9 process, 999 occupy",
		// 	args{
		// 		size:         9,
		// 		occupyNumber: 999,
		// 		resource:     newResource(),
		// 	},
		// },

		// {
		// 	"99 process, 999 occupy",
		// 	args{
		// 		size:         99,
		// 		occupyNumber: 999,
		// 		resource:     newResource(),
		// 	},
		// },

		// {
		// 	"99 process, 9999 occupy",
		// 	args{
		// 		size:         99,
		// 		occupyNumber: 9999,
		// 		resource:     newResource(),
		// 	},
		// },

	}

	for _, tt := range tests {
		t.Logf("运行 %s", tt.name)
		start(tt.args.size, tt.args.occupyNumber, tt.args.resource)
		r := tt.args.resource
		for i := 1; i < tt.args.occupyNumber; i++ {
			if (r.timeOrder[i-1] > r.timeOrder[i]) ||
				(r.timeOrder[i-1] == r.timeOrder[i] && r.processOrder[i-1] > r.processOrder[i]) {
				t.Errorf("%s: resorce 的占用顺序不是按时间排序的", tt.name)
			}
		}
	}
}
