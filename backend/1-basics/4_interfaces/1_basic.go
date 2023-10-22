package main

import (
	"fmt"
)

// структура не знает, какие интерфейсы она реализует
// не нужно указывать в объявлении структуры, какой интерфейс она реализует

type Payer interface {
	Pay(int) error
}

type Wallet struct {
	Cash int
}

func (w *Wallet) Pay(amount int) error {
	if w.Cash < amount {
		return fmt.Errorf("Не хватает денег в кошельке")
	}
	w.Cash -= amount
	return nil
}

func Buy(p Payer) { // вызывается Pay() Wallet'a
	err := p.Pay(10) // имеем доступ только к методам интерфейса Payer(то есть к Cash или к др возможным сущностям конкртной структуры обратиться не можем)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Спасибо за покупку через %T\n\n", p)
}

func main() {
	myWallet := &Wallet{Cash: 100}
	Buy(myWallet)
}
