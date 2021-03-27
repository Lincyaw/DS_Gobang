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

func Seat(clients map[*websocket.Conn]*User, conn *websocket.Conn, msg string) (map[string]string, error) {
	data := map[string]int{"position": 0}
	err := error(nil)
	re := make(map[string]string)
	err = json.Unmarshal([]byte(msg), &data)
	log.Println(data, msg)
	if err != nil {
		log.Println("无法解析的报文：", msg, "error:", err)
		return re, err
	}
	a, b := 0, 0
	for _, v := range clients {
		switch v.Type {
		case 1:
			a = 1
		case 2:
			b = 1
		}
		//if k == conn {
		//	err = errors.New("已经坐下了")
		//	return
		//}
	}
	if a == 1 && b == 1 {
		return re, errors.New("已经坐满了")
	}

	if a == 0 {
		if data["position"] == 1 {
			clients[conn].Type = 1
		} else {
			re["message"] = "这里已经被坐了"
			return re, errors.New("这里已经被坐了")
		}
		a = 1
		re["message"] = clients[conn].Name + "执白子！"
	} else if b == 0 {
		if data["position"] == 2 {
			clients[conn].Type = 2
		} else {
			re["message"] = "这里已经被坐了"
			return re, errors.New("这里已经被坐了")
		}
		b = 1
		re["message"] = clients[conn].Name + "执黑子！"
	}
	if a == 1 && b == 1 { //坐下后坐满了，开局
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
