package main

import (
	"fmt"
	"time"
)

/**

//这是为什么呢，看似合理的程序，是忽略了Channel是阻塞的，如果没有使用go Channel就一直在阻塞的状态，执行就死循环了。这个特性也在很多场合带来了方便。
import "fmt"
func foo(){
defer fmt.Println("World")
fmt.Println("Hello")
}
//func sum(x,y int,c chan int){
//    c <- x + y
//}
func main(){
	foo()
	//    c := make (chan int);
	//    go sum(24,18,c)
	//    fmt.Println(<-c);
	c := make (chan int)
	d := 2
	c <- d+3
	fmt.Println(<-c)
}
**/


//解决死锁的一种方法
//https://medium.com/@zufolo/a-pattern-for-limiting-the-number-of-goroutines-in-execution-56e13b226e72
func main() {

	//channel for terminating the workers
	killq := make(chan bool)

	//queue of jobs
	q := make(chan int)
	// done channel takes the result of the job
	done := make(chan bool)

	numberOfWorkers := 2
	for i := 0; i < numberOfWorkers; i++ {
		go worker(q, i, done, killq)
	}

	numberOfJobs := 17
	for j := 0; j < numberOfJobs; j++ {
		go func(j int) {
			q <- j
		}(j)
	}

	// a deadlock occurs if c >= numberOfJobs
	for c := 0; c < numberOfJobs; c++ {
		<-done
	}

	fmt.Println("finished")

	// cleaning workers
	close(killq)
	time.Sleep(2 * time.Second)
}

func worker(queue chan int, worknumber int, done, ks chan bool) {
	for true {
		select {
		case k := <-queue:
			fmt.Println("doing work!", k, "worknumber", worknumber)
			done <- true
		case <-ks:
			fmt.Println("worker halted, number", worknumber)
			return
		}
	}
}
