package main

import (
	"fmt"
	"time"

	"alanpinder.com/rxgo/v2/rx"
)

func main() {
	x1 := rx.Pipe(rx.Of("hello", "cruel", "world"), rx.ToAny[string]())
	x2 := rx.Pipe(rx.Timer(2*time.Second), rx.ToAny[time.Time]())
	x2 = rx.Pipe(x2, rx.Take[any](5))
	x3 := rx.Pipe(rx.Timer(5*time.Second), rx.ToAny[time.Time]())
	x3 = rx.Pipe(x3, rx.Take[any](2))

	s := rx.NewAutoCompleteSubject[string]()

	x4 := rx.Pipe(s, rx.ToAny[string]())

	s.Signal("howza")

	channel, unsubscribe := rx.CombineLatest(x1, x2, x3, x4).Subscribe()
	defer unsubscribe()

	println("here")

	for x := range channel {
		println(fmt.Sprintf("%v", x))
	}

	println("end")
}

// func main2() {

// 	x := rx.NewAutoCompleteSubject[string]()
// 	c1 := x.Subscribe()
// 	c2 := x.Subscribe()

// 	c3 := x.Subscribe()

// 	xx := sync.WaitGroup{}
// 	xx.Add(1)
// 	xx.Add(1)
// 	xx.Add(1)
// 	xx.Add(1)

// 	go func() {
// 		for x := range c1 {
// 			println(x.String())
// 		}
// 		println("c1 over")
// 		xx.Done()
// 	}()

// 	go func() {
// 		for x := range c2 {
// 			println(x.String())
// 		}
// 		println("c2 over")
// 		xx.Done()
// 	}()

// 	go func() {
// 		for x := range c3 {
// 			println(x.String())
// 		}
// 		println("c3 over")
// 		xx.Done()
// 	}()

// 	x.Signal("hello")

// 	c4 := x.Subscribe(nil)
// 	go func() {
// 		for x := range c4 {
// 			println(x.String())
// 		}
// 		println("c4 over")
// 		xx.Done()
// 	}()

// 	xx.Wait()

// 	// ob := rx.Of(1, 2, 3, 4, 5)
// 	// ob = rx.Filter(func(value int) bool {
// 	// 	return value == 2 || value == 3
// 	// })(ob)
// 	// ob = rx.Map(func(value int) int {
// 	// 	return value * 10
// 	// })(ob)

// 	// ob = rx.ConcatMap(func(value int) rx.Observable[int] {
// 	// 	return rx.Of(value, value+1, value+2)
// 	// })(ob)

// 	// ob = rx.TapDebug[int]("test", fmt.Println)(ob)

// 	// subscription, unsubscribe := ob.Subscribe()

// 	// for {

// 	// 	//		println(fmt.Sprintf("%v\n", rx.Inflight()))

// 	// 	_, open := <-subscription
// 	// 	if !open {
// 	// 		break
// 	// 	}
// 	// 	//		println(fmt.Sprintf("%s: %b\n", value, open))
// 	// }

// 	// unsubscribe()

// 	// //	println(fmt.Sprintf("%v\n", rx.Inflight()))

// 	// time.Sleep(4 * time.Second)

// 	// println(fmt.Sprintf("%v\n", rx.Inflight()))

// }
