package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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
	channel := make(chan def.Context)
	connection, error := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))

	//最後に閉じる
	defer connection.Close()
	defer close(channel)

	if error != nil {
		fmt.Println("接続に失敗しました。")
		panic(error)
	}

	fmt.Printf("接続成功! UUID: %s\n", playerUUID)

	go waitMessage(connection)
	go sendMessage(connection, channel)

	channel <- generateJoin()

	handleInput(channel)

	fmt.Println("Connection lost. exit soon")
	time.Sleep(5 * time.Second)
}

func sendMessage(connection net.Conn, channel <-chan def.Context) {
	context := <-channel

	json := toJSON(context)

	if json != nil {
		_, err := connection.Write(json)

		if err != nil {
			fmt.Println(err)
		}
	}

	sendMessage(connection, channel)
}

func handleInput(channel chan<- def.Context) {
	stdin := bufio.NewScanner(os.Stdin)
	if stdin.Scan() == false {
		return
	}

	channel <- generateContext(def.Message, stdin.Text())

	handleInput(channel)
}

func waitMessage(connection net.Conn) {
	buf := make([]byte, 4*1024)
	n, err := connection.Read(buf)

	if err != nil && err != io.EOF {
		panic(err)
	} else {
		context := new(def.Context)
		if err := json.Unmarshal(buf[:n], context); err != nil {
			fmt.Println("Couldn't unmarshal json.", err)
		} else {
			time, err := time.Parse(time.RFC1123, context.Timestamp)
			formattedTime := time.Format("15:04:05")

			if err != nil {
				fmt.Println("Couldn't parse timestamp:", time)
			} else {
				fmt.Printf("[%s][%s] %s\n", context.UUID, formattedTime, context.Body)
			}
		}

	}
	waitMessage(connection)
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
