package main

import (
	"errors"
	"fmt"
	"time"
)

func readFromCh(ch <-chan int) (int, error) { // receive-only ch
	data, ok := <-ch
	if !ok {
		return 0, errors.New("no data")
	} else {
		return data, nil
	}
}

func writeToTheCh(ch chan<- int) { // send-only ch
	ch <- 2
}

func chanFunc(ch chan int) {}

func TryReceive(c <-chan int) (data int, more, ok bool) {
	select {
	case data, more := <-c:
		return data, more, true
	default:
		return 0, true, false
	}
} // instant return если нет данных

func TryReceiveWithTimeout(ch <-chan int, duration time.Duration) (data int, more, ok bool) {
	select {
	case data, more := <-ch:
		return data, more, true
	case <-time.After(duration): // After() returns a channel that blocks immediately for duration, then returning current time
		return 0, true, false
	}
} // пытаемся получить данные с таймаутом

func Forgotcha() {
	for i := 0; i != 10; i++ {
		// i := i

		go func() {
			fmt.Println(&i, i)
		}()
		// time.Sleep(1 * time.Second) // если не дожидаться горутин внутри или вне цикла, они не успеют выести значение
		// причем если дожидаться их через sleep/waitgroup и тд внутри цикла
		// то каждая горутина успеет вывести конкретное значение локальной перменной i
		// будет 0, 1, ..., 9, адрес перменной не меняется
		// но если никак не дожидаться, то цикл пройдет быстрее, и каждая горутина
		// выведет same address + same last value 9
	}
	time.Sleep(2 * time.Second) // если нигде не ждать горутины, то значения не выведутся
	fmt.Println("end")
}

// fan-out, spreading dataflow from one channel to several
func Fanout(In <-chan int, outA, outB chan int) {
	for data := range In { // receive while channel is opened
		fmt.Println("inside for, data: ", data)
		select {
		case outA <- data:
			fmt.Println("write data to outA")
		case outB <- data:
			fmt.Println("write data to outB")
		}
	}
}

func TestFanOut() { // можем инвокать Fanout до и после акта записи, потому что
	in := make(chan int, 2) // канал buffered --> при записи мы пишем в буфер, если же пытаемся писать в unbuffered канал, то сендер пытается прокинуть данные сразу в стек другой горутины, а если никакая до сих пор не запущена, то планировщик не может отправить ни одну горутину в runnuble state
	outA, outB := make(chan int), make(chan int)
	// var outA_, outB_ chan int - ТАК НИКОГДА НЕ ИНИЦИАЛИЗИРУЕМ, ПОУЛЧАЕМ NIL УКАЗАТЕЛЬ
	// ТК ПАМЯТЬ ДЛЯ КАНАЛА АЛЛОЦИРУЕТСЯ НА КУЧЕ, ВСЕГДА НАДО ЧЕРЕЗ MAKE
	// fmt.Println(outA, outB, outA_, outB_)

	go func(In <-chan int, outA, outB chan int) {
		for data := range In {
			fmt.Println("inside for, data: ", data)
			select {
			case outA <- data:
				fmt.Println("write data to outA")
			case outB <- data:
				fmt.Println("write data to outB")
			}

		}
	}(in, outA, outB)
	// go Fanout(in, outA, outB) // можем инициализировать горутину как до записи в каналы в основном потоке выполнения
	// так и после, причем инициализация и акт записи могут быть отделены в рантайме
	// слипом или другой работой, блокировка будет на время выполнения этой работы
	// но не навсегда, то есть без deadlock'a

	// Fanout(in, ...) когда же пробуем пустить обычную функцию до акта записи
	// будет deadlock, т.к. внутри fanout() пытаемся счиать из канала, в который никто на запишет данные
	time.Sleep(1 * time.Second)
	in <- 1
	in <- 2
	// Fanout(in, ...) когда же пробуем пустить обычную функцию после актов записи, гоуртин нет
	// значит, из канала читаьь будет некому --> deadlock
	fmt.Println("dowork")
	fmt.Println("dowork")
	fmt.Println("dowork")
	fmt.Println("dowork")
	fmt.Println("dowork")
	fmt.Println("dowork")
	fmt.Println("dowork")
	fmt.Println("dowork")
	time.Sleep(1 * time.Second)

	fmt.Println(<-outA, <-outB)
}

func WriteToChanInsideFunc(receiver1, receiver2 chan int) {

	go func() {
		//for data := range receiver1 { // т.к. читаем только из 1-ого ресивера
		//	fmt.Println("got from receiver1", data) // 2-ой в итоге полностью блокирует текущую горутину
		//}
		<-receiver1
		<-receiver2
		fmt.Println("got data from rec")
	}() // хоть какая-то горутина должна быть запущена

	receiver1 <- 1 // момент блокировки
	receiver2 <- 2

	time.Sleep(1 * time.Millisecond)
}

func main() {
	TestFanOut()
	// time.Sleep(1 * time.Second)
	outA, outB := make(chan int), make(chan int)
	WriteToChanInsideFunc(outA, outB)
	// time.Sleep(1 * time.Second)
	// go func() {
	// 	<-outA
	// 	<-outB
	// }()
}
