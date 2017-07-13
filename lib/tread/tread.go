package tread

import (
	"errors"
	"runtime"
)

var CallQueueCap = 16

var (
	callQueue chan func()
	respChan  chan interface{}
)

func checkRun() {
	if callQueue == nil {
		panic(errors.New("mainthread: did not call Run"))
	}
}

func Run(run func()) {
	callQueue = make(chan func(), CallQueueCap)
	respChan = make(chan interface{})

	done := make(chan struct{})
	go func() {
		run()
		done <- struct{}{}
	}()

	for {
		select {
		case f := <-callQueue:
			f()
		case <-done:
			return
		}
	}
}

func CallNonBlock(f func()) {
	checkRun()
	callQueue <- f
}

func Call(f func()) {
	checkRun()
	callQueue <- func() {
		f()
		respChan <- struct{}{}
	}
	<-respChan
}

func CallErr(f func() error) error {
	checkRun()
	callQueue <- func() {
		respChan <- f()
	}
	err := <-respChan
	if err != nil {
		return err.(error)
	}
	return nil
}

func init() {
	runtime.LockOSThread()
}
