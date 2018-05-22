package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

const tasksNum = 20
const workersNum = 3

func cleanUp(tasks chan Work) {

	fmt.Println("clean up!!!")
	close(tasks)

}

func main() {

	wg := new(sync.WaitGroup)
	s := make(chan os.Signal)
	defer close(s)

	signal.Notify(s, os.Interrupt, syscall.SIGTERM)

	tasks := make(chan Work)

	wg.Add(tasksNum)

	go func() {
		for i := 0; i < tasksNum; i++ {
			tasks <- func(wokerId int) {
				fmt.Println("doing task on worker", wokerId)
				time.Sleep(500 * time.Millisecond)
				fmt.Println("finished task on worker", wokerId)
				wg.Done()
			}
		}
		close(tasks)
	}()

	go func() {
		for j := 0; j < workersNum; j++ {
			go worker(j, tasks)
		}
	}()

	go func() {
		<-s
		cleanUp(tasks)
		os.Exit(1)

	}()

	wg.Wait()

	fmt.Println("done!!!")

}
