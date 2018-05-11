package mutual

import (
	"reflect"
	"testing"
)

func Test_start(t *testing.T) {
	type args struct {
		size         int
		occupyNumber int
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
			},
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := start(tt.args.size, tt.args.occupyNumber); !reflect.DeepEqual(got, occupyOrder) {
				t.Errorf("start() = %v, want %v", got, occupyOrder)
			}
		})
	}
}
