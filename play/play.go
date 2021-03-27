package play

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type pData struct {
	Type string `json:"type"`
	pos  `json:"position"`
}
type pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func Play(chessboard [][]int, data string, user int) (map[string]string, error) {
	log.Println(data)
	d := pData{
		Type: "",
		pos:  pos{},
	}
	re:=make(map[string]string)
	err := error(nil)
	err = json.Unmarshal([]byte(data), &d)
	if err != nil {
		log.Println("无法解析的报文：", data, "error:", err)
		return re, err
	}
	fmt.Println("struct", d)

	height, width := len(chessboard), len(chessboard[0])
	if d.X < 0 || d.X >= width || d.Y >= height || d.Y < 0 {
		re["message"] = "不在棋盘范围内"
		return re, err
	}
	if chessboard[d.X][d.Y] == 0 {
		chessboard[d.X][d.Y] = user
	} else {
		re["message"] = "已经在这里下过棋了"
		return re, errors.New("重复下棋")
	}

	// todo: 判断是否有赢家或平手
	for i:=0;i<len(chessboard);i++{
		for j:=0;j<len(chessboard[0]);j++{
			fmt.Print(chessboard[i][j])
		}
		fmt.Println("")
	}
	return re,nil
}
