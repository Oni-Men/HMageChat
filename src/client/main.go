package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"../def"
	"github.com/google/uuid"
)

var (
	playerUUID  string
	displayName string
	port        int
	host        string
)

func processError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	flag.IntVar(&port, "port", 21126, "接続するホストのポート番号")
	flag.StringVar(&host, "host", "localhost", "接続するホストのIPアドレス")
	flag.StringVar(&playerUUID, "uuid", "", "プレイヤーのUUID")
	flag.StringVar(&displayName, "displayname", "", "プレイヤーの名前")
	flag.Parse()

	if playerUUID == "" {
		playerUUID = uuid.New().String()
	}

	fmt.Printf("%s:%dに接続を開始します。\n", host, port)

	//ゴルーチン間でコンテキストをやり取りするためのチャネル
	exit := make(chan bool)
	channel := make(chan def.Context)
	connection, error := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))

	receiver := Receiver{Connection: connection, UUID: playerUUID, DisplayName: displayName}

	//最後に閉じる
	defer connection.Close()
	defer close(channel)
	defer close(exit)

	if error != nil {
		fmt.Println("接続に失敗しました。")
		panic(error)
	}

	fmt.Printf("接続成功! UUID: %s\n", playerUUID)

	go receiver.WaitMessage()
	go sendMessage(connection, channel)

	channel <- generateJoin().Normalize()

	handleInput(channel)

	fmt.Println("Connection lost. exit soon")
	time.Sleep(5 * time.Second)
}

func sendMessage(connection net.Conn, channel <-chan def.Context) {
	for {
		context := <-channel
		json := toJSON(context)

		if json == nil {
			continue
		}

		_, err := connection.Write(json)

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func handleInput(channel chan<- def.Context) {
	stdin := bufio.NewScanner(os.Stdin)
	if stdin.Scan() == false {
		return
	}

	channel <- generateContext(def.Message, stdin.Text())

	handleInput(channel)
}

func generateContext(ctxType def.ContextType, body string) def.Context {
	return def.Context{
		Type:        ctxType,
		UUID:        playerUUID,
		DisplayName: displayName,
		Timestamp:   time.Now().Format(time.RFC1123),
		Body:        body,
	}
}

func toJSON(v interface{}) []byte {
	json, err := json.Marshal(v)

	if err != nil {
		fmt.Println("Couldn't marshal context.", err)
		return nil
	}

	return json
}

func generateJoin() def.Context {
	return generateContext(def.Join, "はまげ")
}
