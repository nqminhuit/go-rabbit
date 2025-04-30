package main

import (
	"fmt"
)

func toChannel(nums *[]int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range *nums {
			fmt.Printf("put %d in\n", n)
			out <- n
		}
		close(out)
	}()
	return out
}

func batch(p <-chan int, batchSize int) <-chan []int {
	out := make(chan []int)
	go func() {
		nums := make([]int, 0, batchSize)
		for n := range p {
			if len(nums) < batchSize {
				nums = append(nums, n)
			} else {
				out <- nums
				nums = make([]int, 0, batchSize)
				nums = append(nums, n)
			}
		}
		out <- nums // this line is important, all leftover elements will not be sent out without this line.
		close(out)
	}()
	return out
}

func batchExample() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 0}

	ch := toChannel(&nums)
	ch2 := batch(ch, 3)

	for x := range ch2 {
		fmt.Printf("\tprocessing %v\n", x)
	}
}
