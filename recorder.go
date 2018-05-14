package mutual

func newRecorder() *[]int {
	res := &[]int{}
	channel := make(chan int, 1000)

	go func() {
		for {
			*res = append(*res, <-channel)
		}
	}()

	return res
}
