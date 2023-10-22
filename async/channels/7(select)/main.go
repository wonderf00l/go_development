package main

import (
	"fmt"
	"time"
)

var start time.Time

func init() {
	start = time.Now()
}

func service1(c chan string) {
	time.Sleep(1 * time.Second)
	c <- "Hello from service 1"
}

func service2(c chan string) {
	time.Sleep(1 * time.Second) // т.к. горутина2 успевает разблокироваться быстрее
	c <- "Hello from service 2" // она быстрее запишет данные в канал -> выполнится ее case
}

func main() {
	fmt.Println("main() started", time.Since(start))
	chan1 := make(chan string)
	chan2 := make(chan string)

	go service1(chan1)
	go service2(chan2)

	// time.Sleep(8 * time.Second) -- если заблокируем главную горутину, в ход пойдут запланированные
	// сразу отпринтят, успеют отоспаться, сразу запишут в канал
	// в select сразу будет выполняться случайным образом case1 || case2

	select {
	case res := <-chan1:
		fmt.Println("Response from service 1", res, time.Since(start))
	case res := <-chan2:
		fmt.Println("Response from service 2", res, time.Since(start))
		/*
				default:
			        fmt.Println("No response received", time.Since(start)) -- default делает блок select не блокирующей операцией, в нашем случае он сразу отработает
					а если бы горутин вообще не было запланировано через go,
					то благодаря default не было бы deadlock
		*/
	}
	/*возможные сценарии:
	1. пробуем прочитать данные из chan1 -> блокируемся ->
	1.1 service1 перехватил инициативу быстрее -> слипом заблокался на 5 сек ->
	-> в этот момент main и sevice1 в блоке -> шедулер толкает service2
	-> он тоже блокается, но на 1 сек, далее блокается при попытке записать в канал
	-> в main блок сразу падает, т.к. срабатывает 2-й case, ведь в chan2 хотят записать данные
	-> и тут уже service2 дорабатывает, main завершает работу
	1.1 service2 оказался быстрее -> блок на 1 сек -> падаем в service1 -> блок на 5 сек
	-> service2 размораживается -> но блочится при попытке записать не в ch1
	тут уже main разблочиться с выбранным case2 -> исход тот же
	2. пробуем в chan2
	2.1 s1 быстрее -> блок -> далее блок s2 но проходит быстрее -> запись в родной канал
	2.2. s2 быстрее -> блок 1 сек -> блок в s1 -> s2 просыпается и пишет в родной канал - НАИБОЛЕЕ БЫСТРЫЙ ИСХОД
	*/
	fmt.Println("main() stopped", time.Since(start))
} // при одинаковом слипе -> исход не определен

/*
Чтение из буферизированного канала не является
блокируемой операцией, если буфер не пустой,
поэтому все блоки case будут неблокируемыми,
и во время выполнения Go выберет case случайным образом.
*/

// пустой select блокирует горутину навсегда
