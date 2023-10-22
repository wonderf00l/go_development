package main

import (
	"fmt"
	"unicode/utf8"
	"golang.org/x/text/unicode/norm"
)

// под капотом строки - слайсы байт

func main() {
	// пустая строка по-умолчанию
	var str string
	norm.String()

	// со спец символами
	var hello string = "Привет\n\t"

	// без спец символов, переносы отображаются
	var world string = `Мир\n\t 
	jtjtjrejlreer
	ertert
	ert
	er
	tsdsdff
	dfef
	`

	fmt.Println("str", str)
	fmt.Println("hello", hello)
	fmt.Println("world", world)

	// UTF-8 из коробки
	var helloWorld = "Привет, Мир!"
	hi := "你好，世界 <- 💩"

	fmt.Println("helloWorld", helloWorld)
	fmt.Println("hi", hi)

	// одинарные кавычки для байт (uint8): utf-7/8
	var rawBinary byte = '\x27'
	fmt.Printf("raw binary: %b", rawBinary)

	// строки к слайсу byte/rune

	// rune (uint32) для utf-7/UTF-8/utf-16/utf-32 символов(1-, 2-, 4-байтовые строки)
	var someChinese rune = '茶'

	fmt.Println("info: ", rawBinary, someChinese)

	helloWorld = "Привет Мир 👋ё"
	// конкатенация строк
	andGoodMorning := helloWorld + " и доброе утро!"

	fmt.Println(helloWorld, andGoodMorning)

	// строки неизменяемы
	// cannot assign to helloWorld[0]
	// helloWorld[0] = 72

	// получение длины строки
	fmt.Println("get len in bytes/runes")
	fmt.Println(utf8.RuneCountInString("string"), '\n') // 6 10, 6 - длина строки --> используем rune функции для работы со строками в общем случае
	byteLen := len(helloWorld)                          // 26 байт
	symbols := utf8.RuneCountInString(helloWorld)       // 13 рун

	fmt.Println(byteLen, symbols)

	// получение подстроки, в байтах, не символах!
	hello = helloWorld[:12] // Привет, 0-11 байты, копируются БАЙТЫ, НЕ СИМВОЛЫ!
	fmt.Println("hello:", hello, '\n')
	H := helloWorld[0] // byte, 208, не "П"(код 'П' - 1055, то есть отобразилась только часть, вмещающаяся в байт)
	fmt.Println("first ch:", H)
	fmt.Println("convert to []rune, now 1st ch code:", []rune(helloWorld)[0])      // 1055
	fmt.Printf("convert to []rune, now 1st ch as sym:%c\n", []rune(helloWorld)[0]) // П
	for _, sym := range helloWorld {
		fmt.Printf("%T\n", sym) // int32
	}
	for _, sym := range "string" {
		fmt.Printf("%T\n", sym) // int32
	}

	// конвертация в слайс байт и обратно
	byteString := []byte(helloWorld) // ВЗАИМНАЯ КОНВЕРТАЦИЯ
	helloWorld = string(byteString)

	fmt.Println(byteString, helloWorld)

	for i, ch := range "string" {
		fmt.Printf("%d %c\n", i, ch)
	}
}
