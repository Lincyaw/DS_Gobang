package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"play/play"
	"strconv"
)

var Chessboard = make([][][]int, 25) //房间 列 行
var ChessInUse = make([]bool, 25)
var clients = make(map[*websocket.Conn]*play.User) // 客户端，标识：1,2 玩家，3观众，0未连接
var broadcast = make(chan play.Message)            // broadcast channel
var state = make([]chan int, 25)

// 只配置了跨域请求，其他配置见文档
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	//创建棋盘
	for k, _ := range Chessboard {
		Chessboard[k] = make([][]int, 15)
		for j := range Chessboard[k] {
			Chessboard[k][j] = make([]int, 15)
		}
	}
	for k, _ := range state {
		state[k] = make(chan int, 2)
		state[k] <- 1
	}
	// 创建静态资源服务
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// 创建websocket路由
	http.HandleFunc("/ws", handleConnections)

	// 开始监听
	go handleMessages()

	// 开启服务端
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//ws处理器
func handleConnections(w http.ResponseWriter, r *http.Request) {
	var roomNumber int = 0
	//升级get请求到ws请求
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()
	// todo: 连接关闭后，标志该房间被释放

	// 注册客户端
	clients[ws] = &play.User{Name: "用户" + strconv.Itoa(int(rand.Int31n(1000000))), Type: 3}

	for {
		// Read in a new message as JSON and map it to a Message object
		messageType, msg, err := ws.ReadMessage()
		log.Println("收到消息", string(msg), messageType)
		if err != nil {
			log.Printf("error: %v - %d", err, messageType)
			delete(clients, ws)
			break
		}
		//只接受文字消息
		if messageType != 1 {
			continue
		}
		data := make([]string, 2, 2)
		for k, v := range msg {
			if string(v) == "|" {
				data[0] = string(msg[0:k])
				data[1] = string(msg[k+1:])
			}
		}
		log.Println("解析消息", data)
		reData := new(play.Message)
		reData.User = *clients[ws]

		switch data[0] {
		// 清空棋盘
		case "clear":
			for i := range Chessboard[roomNumber] {
				for j := range Chessboard[roomNumber][i] {
					Chessboard[roomNumber][i][j] = 0
				}
			}
			for i := 0; i < len(Chessboard[roomNumber]); i++ {
				for j := 0; j < len(Chessboard[roomNumber][i]); j++ {
					fmt.Print(Chessboard[roomNumber][i][j])
				}
				fmt.Println("")
			}

		case "message":
			reData.Type = "message"
			reData.Data, err = play.SendMessage(data[1])
		case "play":
			log.Print("以下房间已经被使用：")
			for i := range ChessInUse {
				if ChessInUse[i] {
					log.Print(i)
				}
			}
			log.Print("\n")
			reData.Type = "play"
			if clients[ws].Type != 3 {
				log.Println("state before", state[roomNumber])
				log.Println("roomNum", roomNumber)
				// 获取状态
				s := <-state[roomNumber]
				switch clients[ws].Type {
				case 1:
					if s == 1 {
						reData.Data, err = play.Play(Chessboard[roomNumber], data[1], clients[ws].Type)
						if err != nil {
							state[roomNumber] <- 1
						} else {
							state[roomNumber] <- 0
						}
					} else {
						state[roomNumber] <- s
						reData.Data = map[string]string{"message": "不是你下的时候"}
					}
				case 2:
					if s == 0 {
						reData.Data, err = play.Play(Chessboard[roomNumber], data[1], clients[ws].Type)
						// 如果有错误，就还是本方下棋
						if err != nil {
							state[roomNumber] <- 0
						} else {
							state[roomNumber] <- 1
						}
					} else {
						state[roomNumber] <- s
						reData.Data = map[string]string{"message": "不是你下的时候"}
					}
				}
				log.Println("state after", state[roomNumber])
			}
		case "seat":
			reData.Type = "radio"
			log.Println("seat", roomNumber)
			reData.Data, err = play.Seat(clients, ws, data[1], roomNumber)
		case "changeName":
			reData.Type = "radio"
			reData.Data, err = play.ChangeName(clients[ws], data[1])
		case "chooseRoom":
			roomNumber, err = strconv.Atoi(data[1])
			if err == nil {
				ChessInUse[roomNumber] = true
			} else {
				for i := range ChessInUse {
					if !ChessInUse[i] {
						roomNumber = i
						ChessInUse[i] = true
						break
					}
				}
			}
			clients[ws].Room = roomNumber

		default:
			continue
		}
		if err != nil {
			log.Printf("任务异常，忽略: %v ", err)
			continue
		}
		log.Println(reData)

		// 发送到通道
		broadcast <- *reData
	}
}

//监听消息
func handleMessages() {
	for {
		// 接受消息
		data := <-broadcast
		// 广播消息
		for client, user := range clients {
			if user.Room != data.Room {
				continue
			}
			re, _ := json.Marshal(data)
			//log.Println("广播：",re)
			err := client.WriteMessage(websocket.TextMessage, re)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
