package play

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

var MsgNewName = map[string]string{
	"new": "",
}

func Seat(clients map[*websocket.Conn]*User, conn *websocket.Conn, msg string, roomNum int) (map[string]string, error) {
	data := map[string]int{"position": 0}
	err := error(nil)
	re := make(map[string]string)
	err = json.Unmarshal([]byte(msg), &data)
	log.Println(data, msg, roomNum)
	if err != nil {
		log.Println("无法解析的报文：", msg, "error:", err)
		return re, err
	}
	a, b := [25]int{}, [25]int{}
	for k, _ := range a {
		a[k] = 0
		b[k] = 0
	}
	for k, v := range clients {
		if k == conn {
			if v.Type != 3 {
				re["message"] = "你已经坐下了"
				return re, err
			}
			switch v.Type {
			case 1:
				a[roomNum] = 1
			case 2:
				b[roomNum] = 1
			}
		}

	}

	if a[roomNum] == 0 && b[roomNum] == 0 {
		clients[conn].Type = data["position"]
		re["message"] = clients[conn].Name + "进入房间！"
	} else if a[roomNum] == 0 {
		if data["position"] == 1 {
			clients[conn].Type = 1
			a[roomNum] = 1
			re["message"] = clients[conn].Name + "执白子！"
		} else {
			re["message"] = "这里已经被坐了"
			return re, errors.New("这里已经被坐了")
		}
	} else if b[roomNum] == 0 {
		if data["position"] == 2 {
			clients[conn].Type = 2
			b[roomNum] = 1
			re["message"] = clients[conn].Name + "执黑子！"
		} else {
			re["message"] = "这里已经被坐了"
			return re, errors.New("这里已经被坐了")
		}
	} else {
		return re, errors.New("已经坐满了")
	}
	if a[roomNum] == 1 && b[roomNum] == 1 { //坐下后坐满了，开局
		//begin()
		fmt.Println("开局")
	}
	return re, nil
}
func ChangeName(user *User, msg string) (re map[string]string, err error) {
	data := map[string]string{"newName": ""}
	err = json.Unmarshal([]byte(msg), &data)
	log.Println(data, msg)
	if err != nil {
		log.Println("无法解析的报文：", msg, "error:", err)
		return
	}
	if data["newName"] == "" {
		err = errors.New("新名字不能为空")
		return
	}
	re = make(map[string]string)
	re["message"] = user.Name + "改名为" + data["newName"] + "了！"
	user.Name = data["newName"]
	return
}
