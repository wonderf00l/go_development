package main

import (
	"fmt"
	"strconv"
)

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

func (w *Wallet) String() string {
	return "Кошелёк в котором " + strconv.Itoa(w.Cash) + " денег"
}

// -----

func main() {
	myWallet := &Wallet{Cash: 100}
	fmt.Printf("Raw payment : %#v\n", myWallet) // go - представление структуры(go - код) #v &main.Wallet{Cash:100}
	fmt.Printf("Способ оплаты: %s\n", myWallet) // cast к стрингеру и вызов string()
	// sringer, ok := myWaller.(fmt.Stringer)
	// if ok {
	// 	stringer.String()
	// }

	var x interface{} = "foo"

	switch v := x.(type) { // v - stores value of x
	case nil:
		fmt.Println("x is nil") // here v has type interface{}
	case int:
		fmt.Println("x is", v) // here v has type int
	case bool, string:
		fmt.Println("x is bool or string: ", v) // here v has type interface{}
	default:
		fmt.Println("type unknown") // here v has type interface{}
	}
}
