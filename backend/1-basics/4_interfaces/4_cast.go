package main

import (
	"fmt"
)

// --------------

type Wallet struct {
	Cash int
}

func (w *Wallet) Pay(amount int) error {
	if w.Cash < amount {
		return fmt.Errorf("Not enough cash")
	}
	w.Cash -= amount
	return nil
}

// --------------

type Card struct {
	Balance    int
	ValidUntil string
	Cardholder string
	CVV        string
	Number     string
}

func (c *Card) Pay(amount int) error {
	if c.Balance < amount {
		return fmt.Errorf("Not enough money on balance")
	}
	c.Balance -= amount
	return nil
}

// --------------

type ApplePay struct {
	Money   int
	AppleID string
}

func (a *ApplePay) Pay(amount int) error {
	if a.Money < amount {
		return fmt.Errorf("Not enough money on account")
	}
	a.Money -= amount
	return nil
}

// --------------

type Payer interface {
	Pay(int) error
}

// --------------

// func Buy(p *Payer) { если требовать указатель на интерфейс в аругментах --> при передаче &Waller{...} ошибка: *Wallet != *Payer
// 	switch (*p).(type) {
// 	case *Wallet:
// 		fmt.Println("Оплата наличными?")
// 	case *Card:
// 		plasticCard, ok := (*p).(*Card) // каст: от интерфейса --> к его реализации, чтобы иметь доступ к содержимому конкретной реализации
// 		if !ok {                        // обычно докапывание до деталей конкретной имплементации интерфейса - bad practice
// 			// не гибко
// 			fmt.Println("Не удалось преобразовать к типу *Card")
// 		}
// 		fmt.Println("Вставляйте карту,", plasticCard.Cardholder)
// 	default:
// 		fmt.Println("Что-то новое!")
// 	}

// 	err := (*p).Pay(10)
// 	if err != nil {
// 		fmt.Printf("Ошибка при оплате %T: %v\n\n", (*p), err)
// 		return
// 	}
// 	fmt.Printf("Спасибо за покупку через %T\n\n", p)
// }

// However, there is a limitation on interfaces.
// If you pass a structure directly into an interface, only value methods of that type (ie. func (f Foo) Dummy(),
// not func (f *Foo) Dummy()) can be used to fulfill the interface. This is because you're storing a copy of the original structure in the
// interface, so pointer methods would have unexpected effects (ie. unable to alter the original structure). Thus the default rule of thumb is
// to store pointers to structures in interfaces, unless there's a compelling reason not to.

// Values wrapped in interfaces are not addressable. When an interface value is created, the value that is wrapped in the interface is copied.
// It is therefore not possible to take its address. Theoretically you could allow to take the address of the copy, but that would be the source
// of (even) more confusion than what benefit it would provide, as the address would point to a copy, and methods with pointer receiver could only
// modify the copy and not the original.

func Buy(p Payer) { // но при требовании просто payer все валидно даже при передаче указателя на структуры-имплементации интерфейса Payer
	switch p.(type) { // в "интерфейс" можем передавать как копию ориг структуры, так и указатель на нее, в зависимости от этого могут быть использованы методы с разными ресиверами
	case *Wallet:
		fmt.Println("Оплата наличными?")
	case *Card:
		plasticCard, ok := p.(*Card) // каст: от интерфейса --> к его реализации, чтобы иметь доступ к содержимому конкретной реализации
		if !ok {                     // обычно докапывание до деталей конкретной имплементации интерфейса - bad practice
			// не гибко
			fmt.Println("Не удалось преобразовать к типу *Card")
		}
		fmt.Println("Вставляйте карту,", plasticCard.Cardholder)
	default:
		fmt.Println("Что-то новое!")
	}

	err := p.Pay(10)
	if err != nil {
		fmt.Printf("Ошибка при оплате %T: %v\n\n", p, err)
		return
	}
	fmt.Printf("Спасибо за покупку через %T\n\n", p)
}

// --------------

func main() {

	myWallet := &Wallet{Cash: 100}
	Buy(myWallet)

	var myMoney Payer
	myMoney = &Card{Balance: 100, Cardholder: "rvasily"}
	Buy(myMoney)

	myMoney = &ApplePay{Money: 9}
	Buy(myMoney)
}
