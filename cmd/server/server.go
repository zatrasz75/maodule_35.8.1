package main

import (
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	proverbs = "https://go-proverbs.github.io/"
	addr     = "localhost:12345"
	proto    = "tcp4"
)

func main() {
	str, err := getStr(proverbs)
	if err != nil {
		log.Fatal(err)
	}
	// Запуск сетевой службы по протоколу TCP
	listener, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			return
		}
	}(listener)

	// Подключения обрабатываются в бесконечном цикле.
	// Иначе после обслуживания первого подключения сервер
	// завершит работу.
	for {
		// Принимаем подключение.
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("установлено соединение:", conn.RemoteAddr())

		// Вызов обработчика подключения.
		go handleConn(conn, str)
	}
}

// Обработчик. Вызывается для каждого соединения.
func handleConn(conn net.Conn, str []string) {
	// Закрытие соединения.
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	// Читаем рандомно строки раз в 3 секунды.
	for {
		_, err := conn.Write([]byte(str[rand.Intn(len(str))] + "\n\r"))
		if err != nil {
			return
		}

		time.Sleep(3 * time.Second)
	}
}

// Получает данные по заданному URL
func getStr(url string) ([]string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := parseUrl(string(body))
	return data, nil
}

// Выбирает нужные строки с Url
func parseUrl(body string) (data []string) {

	tkn := html.NewTokenizer(strings.NewReader(body))
	var values []string
	var isH3, isProverb bool

	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			return values
		case tt == html.StartTagToken:
			t := tkn.Token()
			if !isH3 {
				isH3 = t.Data == "h3"
			} else {
				isProverb = t.Data == "a"
			}
		case tt == html.TextToken:
			t := tkn.Token()
			if isProverb {
				values = append(values, t.Data)
			}
			isH3 = false
			isProverb = false
		}
	}
}
