package main

import (
	
	"fmt"
	"sync"
	"time"
	//"sort"
)

type Work func(workerId int)

func worker(workerID int, tasks chan Work) {

	for {
		select {

		case task, ok := <-tasks:
			if !ok { //signal the workers to shutdown
				fmt.Println(ok)
				return
			}
			task(workerID)
		default:
			fmt.Println("waiting for task...")
			time.Sleep(30 * time.Millisecond)
		}
	}
}

func main() {

	wg := new(sync.WaitGroup)

	tasks := make(chan Work)
	const tasksNum = 10

	wg.Add(tasksNum)

	go func() {
		for i := 0; i < tasksNum; i++ {
			tasks <- func(wokerId int) {
				time.Sleep(500 * time.Millisecond)
				fmt.Println("finished job on worker", wokerId)
				wg.Done()
			}
		}
		close(tasks)
	}()

	go func() {
		for j := 0; j < 3; j++ {
			go worker(j, tasks)
		}
	}()

	wg.Wait()

	
	fmt.Println("done!!!")

}
