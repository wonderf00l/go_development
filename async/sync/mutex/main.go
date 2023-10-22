package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Inner struct {
	a, b int
}

func (in Inner) DoWork() {
	fmt.Println("Inside Inner")
}

type InnerEmbedded struct {
	c, d int
}

func (in InnerEmbedded) DoEmbeddedWork() {
	fmt.Println("inside InnerEmbedded")
}

type Outer struct {
	Inner Inner
	InnerEmbedded
}

type playerWallet struct {
	coins int64
	mu    sync.RWMutex
}

func (w *playerWallet) getCoins() int64 {
	w.mu.RLock()
	coins := w.coins
	w.mu.RUnlock()

	return coins
}

// разделяем общую очередь за ресурсом на очередь на чтение/запись - оптимизация
func (w *playerWallet) spendCoins(amount int64) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	if w.coins-amount < 0 {
		return errors.New("insufficientFounds")
	}

	w.coins -= amount
	log.Printf("spent %d coin(s), balance: %d", amount, w.coins)
	return nil
}

func main() {

	outer := Outer{
		Inner: Inner{
			a: 1,
			b: 2,
		},
		InnerEmbedded: InnerEmbedded{
			c: 3,
			d: 4,
		},
	}

	fmt.Println(outer)

	fmt.Println(outer.Inner.a)
	fmt.Println(outer.c)

	outer.DoEmbeddedWork()
	outer.Inner.DoWork()

	var (
		wg     sync.WaitGroup
		wallet = &playerWallet{coins: 5}
	)

	rand.Seed(time.Now().Unix())

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		payForUnitsAndBuildings(wallet)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		buyAnotherUnitInShop(wallet)
	}()
	// каждый раз "не успевает" совершить валидную операцию разная горутина
	wg.Wait()
}

func payForUnitsAndBuildings(w *playerWallet) {
	if err := w.spendCoins(3); err != nil {
		log.Println(err)
	}
}

func buyAnotherUnitInShop(w *playerWallet) {
	if err := w.spendCoins(4); err != nil {
		log.Println(err)
	}
}

var errInsufficientFunds = errors.New("insufficient funds")

type playerWallet struct {
	coins int64
	mu    sync.Mutex
}

func (w *playerWallet) spendCoins(amount int64) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.coins-amount < 0 {
		return errInsufficientFunds
	}

	w.coins -= amount
	log.Printf("spent %d coin(s), balance: %d", amount, w.coins)
	return nil
}
