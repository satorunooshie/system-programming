package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	fmt.Println("start sub()")
	done := make(chan bool)
	go func() {
		fmt.Println("sub() is finished")
		done <- true
	}()
	go sub()
	time.Sleep(2 * time.Second)

	// 値が来るたびにforループが回る個数が未定の動的配列
	pn := primeNumber()
	for n := range pn {
		fmt.Println(n)
	}
}

func sub() {
	fmt.Println("sub() is running")
	time.Sleep(time.Second)
	fmt.Println("sub() is finished")
}

func primeNumber() (result chan int) {
	result = make(chan int)
	go func() {
		result <- 2
		for i := 3; i < 100000; i += 2 {
			l := int(math.Sqrt(float64(i)))
			hasFound := false
			for j := 3; j < l; j += 2 {
				if i % j == 0 {
					hasFound = true
					break
				}
			}
			if !hasFound {
				result <- i
			}
		}
		close(result)
	}()
	return
}