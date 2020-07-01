package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"./chat"
)

// ServerConfig サーバーの構成
type ServerConfig struct {
	Port int
}

func main() {

	serverConfig := loadServerConfig()

	tcpAddr := resolveAddress(serverConfig.Port)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	handleError(err)

	fmt.Printf("ポート「%d」で接続を待機中です\n", tcpAddr.Port)

	channel := make(chan chat.Notification)
	observer := chat.Observer{Senders: make([]chat.Sender, 0, 5), Subject: channel}

	go observer.WaitNotice()
	waitClient(listener, observer, channel)
}

func waitClient(listener net.Listener, observer chat.Observer, channel chan chat.Notification) {
	connection, error := listener.Accept()
	handleError(error)

	receiver := chat.Receiver{UUID: "", Connection: connection, Observer: channel}
	go receiver.WaitMessage()

	waitClient(listener, observer, channel)
}

func loadServerConfig() ServerConfig {
	serverConfig := ServerConfig{
		Port: 21126,
	}

	filePath := "./server_config.json"

	if !isExist(filePath) {
		_, err := os.Create(filePath)

		if err != nil {
			panic(err)
		}

		buf, err := json.Marshal(serverConfig)

		if err != nil {
			panic(err)
		}

		ioutil.WriteFile(filePath, buf, 777)
	}

	buf, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Println("server_config.jsonを読み込めませんでした")
		panic(err)
	}

	if err = json.Unmarshal(buf, &serverConfig); err != nil {
		fmt.Println("server_config.jsonを解析出来ませんでした")
		panic(err)
	}

	return serverConfig
}

func resolveAddress(port int) *net.TCPAddr {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	handleError(err)
	return address
}

func handleError(e error) {
	if e != nil {
		panic(e)
	}
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
