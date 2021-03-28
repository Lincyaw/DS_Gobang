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
	re := make(map[string]string)
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
		position := pos {
		    X: d.X,
		    Y: d.Y,
		}
		tmp, _ := json.Marshal(position)
		re["position"] = string(tmp)
	} else {
		re["message"] = "已经在这里下过棋了"
		return re, errors.New("重复下棋")
	}

	if checkWin(chessboard, d.X, d.Y) {
		re["message"] = "您获胜了！"
	}
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			fmt.Print(chessboard[i][j])
		}
		fmt.Println("")
	}

	return re, nil
}

func checkByStep(chessboard [][]int, x, y, xdir, ydir int) bool {
	height, width := len(chessboard), len(chessboard[0])
	cnt := 0
	for i := 1; i < 5; i++ {
		tmpX, tmpY := x-xdir*i, y-ydir*i
		if tmpX < 0 || tmpX >= width || tmpY >= height || tmpY < 0 {
			break
		}
		if chessboard[tmpX][tmpY] == chessboard[x][y] {
			cnt++
		}
	}

	for i := 1; i < 5; i++ {
		tmpX, tmpY := x+xdir*i, y+ydir*i
		if tmpX < 0 || tmpX >= width || tmpY >= height || tmpY < 0 {
			break
		}
		if chessboard[tmpX][tmpY] == chessboard[x][y] {
			cnt++
		}
	}
	if cnt >= 4 {
		return true
	}
	return false
}

func checkWin(chessboard [][]int, x, y int) bool {
	if checkByStep(chessboard, x, y, 0, 1) {
		return true
	}
	if checkByStep(chessboard, x, y, 1, 0) {
		return true
	}
	if checkByStep(chessboard, x, y, 1, 1) {
		return true
	}
	if checkByStep(chessboard, x, y, -1, 1) {
		return true
	}
	return false
}

