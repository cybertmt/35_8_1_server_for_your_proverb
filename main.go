package main

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

// Сетевой адрес.
//
// Служба будет слушать запросы на всех IP-адресах
// компьютера на порту 12345.
// Например, 127.0.0.1:12345
const addr = "127.0.0.1:12345"

// Протокол сетевой службы.
const proto = "tcp4"

// Массив пословиц
// и файл с пословицами, записанными построчно
var Verbs []string
const verbsFile = "verbs.txt"

func main() {

	file, err := os.Open(verbsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Verbs = append(Verbs, scanner.Text())
	}

	// Запуск сетевой службы по протоколу TCP
	// на порту 12345.
	listener, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatal(err)
	}
	// Подключения обрабатываются в бесконечном цикле.
	// Иначе после обслуживания первого подключения сервер
	//завершит работу.
	for {
		// Принимаем подключение.
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Вызов обработчика подключения. -- go для многопоточности
		go handleConn(conn)
	}
}

// Обработчик. Вызывается для каждого соединения.
func handleConn(conn net.Conn) {
	go verbs(conn)
	// Чтение сообщения от клиента.
	for {
		reader := bufio.NewReader(conn)
		b, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println(err)
			return
		}
		// Удаление символов конца строки.
		msg := strings.TrimSuffix(string(b), "\n")
		msg = strings.TrimSuffix(msg, "\r")
		// Если получили "q" - закрываем соединение.
		if msg == "q" {
			conn.Write([]byte("Bye\n\r"))
			conn.Close()
		}
		conn.Write([]byte("q to quit\n\r"))
	}
}

func verbs(conn net.Conn) {
	rand.Seed(time.Now().Unix())
	for {
		time.Sleep(3*time.Second)
		conn.Write([]byte(Verbs[rand.Intn(len(Verbs))] + "\n\r"))
	}
}