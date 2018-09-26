package mutual

// func Test_start(t *testing.T) {
// 	type args struct {
// 		processes int
// 		occupieds int
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{

// 		// {
// 		// 	"3 process, 9 occupy",
// 		// 	args{
// 		// 		processes: 3,
// 		// 		occupieds: 9,
// 		// 	},
// 		// },

// 		// {
// 		// 	"20 process, 200 occupy",
// 		// 	args{
// 		// 		processes: 20,
// 		// 		occupieds: 200,
// 		// 	},
// 		// },

// 		// {
// 		// 	"20 process, 400 occupy",
// 		// 	args{
// 		// 		processes: 20,
// 		// 		occupieds: 400,
// 		// 	},
// 		// },

// 		// {
// 		// 	"20 process, 600 occupy",
// 		// 	args{
// 		// 		processes: 20,
// 		// 		occupieds: 600,
// 		// 	},
// 		// },

// 		{
// 			"99 process, 9999 occupy",
// 			args{
// 				processes: 99,
// 				occupieds: 9999,
// 			},
// 		},

// 		//
// 	}

// 	for _, tt := range tests {
// 		t.Logf("运行 %s", tt.name)
// 		r := start(tt.args.processes, tt.args.occupieds)
// 		if len(r.timeOrder) != tt.args.occupieds {
// 			t.Errorf("占用资源的次数不对 occupieds=%d, 实际值=%d", tt.args.occupieds, len(r.timeOrder))
// 		}
// 		for i := 1; i < tt.args.occupieds; i++ {
// 			if (r.timeOrder[i-1] > r.timeOrder[i]) ||
// 				(r.timeOrder[i-1] == r.timeOrder[i] && r.procOrder[i-1] > r.procOrder[i]) {
// 				t.Errorf("%s: resource 的占用顺序不是按时间排序的", tt.name)
// 			}
// 		}
// 	}
// }
