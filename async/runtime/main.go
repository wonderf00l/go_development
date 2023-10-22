package main

import (
	"fmt"
	"runtime"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		// runtime.Gosched()
		/*
			yield new processor and allows other goroutines to run
			фактически говорим, что нужно переключить контекст и дать порабоать
			другим горутинам
			возможно, полезно, если по сути горутина делает не
			блокирующую по своей природе работу(зависла в цикле for например)
			но занимает много процессорного времени
			тем самым врчную делаем context switching
			на blocking operations - automatically(gopark() - goready() for another goroutine)
		*/
		fmt.Println(s)
	}
}

func main() {
	runtime.GOMAXPROCS(1) // default is runtime.NumCPU()
	/*
		When this variable is set to a positive number N, Go runtime will be able to create up
		to N native threads, on which all green threads will be scheduled. Native thread a kind
		of thread which is created by the operating system (Windows threads, pthreads etc).
		This means that if N is greater than 1, it is possible that goroutines will be scheduled
		to execute in different native threads and, consequently, run in parallel
		(at least, up to your computer capabilities: if your system is based on multicore processor,
		it is likely that these threads will be truly parallel; if your processor has single core,
		then preemptive multitasking implemented in OS threads will create a visibility of parallel execution).
	*/
	/*
		при либом значении GOMAXPROCS имеем недетерминированное
		поведение программы: когда то доп горутина успевает отработать(причем всегда полностью)
		когда-то - нет
		Looks like that in newer versions of Go compiler Go runtime forces goroutines to yield
		not only on concurrency primitives usage, but on OS system calls too.
		This means that execution context can be switched between goroutines also on IO functions calls.
		Consequently, in recent Go compilers it is possible to observe indeterministic behavior even when GOMAXPROCS is unset or set to 1.

		если свитчим котнекст с goshed, то принтинг будем поочередным

		фактически, значение GOMAXPROCS влияет лишь на то,
		будет ли это concurrency в рамках 1-ого OS потока,
		либо же нескольких(и если multi-core processor, следовательно, будет ли параллелизм)

		https://blog.ildarkarymov.ru/posts/go-concurrency/#%D0%BA%D0%B0%D0%BA-%D0%BF%D0%BB%D0%B0%D0%BD%D0%B8%D1%80%D0%BE%D0%B2%D1%89%D0%B8%D0%BA-go-%D1%83%D0%BF%D1%80%D0%B0%D0%B2%D0%BB%D1%8F%D0%B5%D1%82-%D0%B3%D0%BE%D1%80%D1%83%D1%82%D0%B8%D0%BD%D0%B0%D0%BC%D0%B8
	*/
	go say("world")
	say("hello")
}
