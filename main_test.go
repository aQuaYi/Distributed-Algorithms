package mutual

import (
	"sort"
	"testing"
)

func Test_start(t *testing.T) {
	type args struct {
		size         int
		occupyNumber int
		rsc          *resource
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
				rsc:          newResource(),
			},
		},

		{
			"9 process, 99 occupy",
			args{
				size:         9,
				occupyNumber: 99,
				rsc:          newResource(),
			},
		},

		{
			"3 process, 999 occupy",
			args{
				size:         3,
				occupyNumber: 999,
				rsc:          newResource(),
			},
		},

		{
			"6 process, 999 occupy",
			args{
				size:         6,
				occupyNumber: 999,
				rsc:          newResource(),
			},
		},

		{
			"9 process, 999 occupy",
			args{
				size:         9,
				occupyNumber: 999,
				rsc:          newResource(),
			},
		},

		{
			"99 process, 999 occupy",
			args{
				size:         99,
				occupyNumber: 999,
				rsc:          newResource(),
			},
		},

		{
			"3 process, 99 occupy",
			args{
				size:         3,
				occupyNumber: 99,
				rsc:          newResource(),
			},
		},

		{
			"5 process, 500 occupy",
			args{
				size:         5,
				occupyNumber: 500,
				rsc:          newResource(),
			},
		},

		{
			"3 process, 999 occupy",
			args{
				size:         3,
				occupyNumber: 999,
				rsc:          newResource(),
			},
		},

		// 	{
		// 		"99 process, 9999 occupy",
		// 		args{
		// 			size:         99,
		// 			occupyNumber: 9999,
		// 			rsc:          newResource(),
		// 		},
		// 	},

	}

	for _, tt := range tests {
		t.Log("运行", tt.name)
		start(tt.args.size, tt.args.occupyNumber, tt.args.rsc)
		if !sort.IntsAreSorted(tt.args.rsc.timeOrder) {
			t.Errorf("%s: resorce 的占用顺序不是按时间排序的", tt.name)
		}
	}
}
