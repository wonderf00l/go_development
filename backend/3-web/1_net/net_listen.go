package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	name := conn.RemoteAddr().String()

	fmt.Printf("%+v connected\n", name)
	conn.Write([]byte("Hello, " + name + "\n\r"))

	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() { // блокаемся пока клиент подключен к нам
		text := scanner.Text()
		if text == "Exit" {
			conn.Write([]byte("Bye\n\r"))
			fmt.Println(name, "disconnected")
			break
		} else if text != "" {
			fmt.Println(name, "enters", text)
			conn.Write([]byte("You enter " + text + "\n\r"))
		}
	}
}

// nc localhost 8080
func main() {
	listner, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listner.Accept() // в синхронном режиме можем работать только с 1-им клиентом, если он отключистя, на этом этапе ждем подлкючение
		if err != nil {
			panic(err)
		}
		go handleConnection(conn) // при запуске асинхронно горутин запускается N=числу подключившихся  клиентов
	}
}
