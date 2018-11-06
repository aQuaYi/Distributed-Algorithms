package main

import (
	"reflect"
	"testing"
)

func TestInt64ToHex(t *testing.T) {
	type args struct {
		num int64
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{

		{
			"24",
			args{
				num: 24,
			},
			[]byte{0, 0, 0, 0, 0, 0, 0, 24},
		},

		{
			"1024",
			args{
				num: 1024,
			},
			[]byte{0, 0, 0, 0, 0, 0, 4, 0},
		},

		//
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64ToHex(tt.args.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64ToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}
