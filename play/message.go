package play

import (
	"encoding/json"
	"log"
)

func SendMessage(msg string) (re map[string]string, err error)  {
	data := map[string]string{"message":""}
	err = json.Unmarshal([]byte(msg),&data)
	if err != nil{
		log.Println("无法解析的报文：",msg,"error:",err)
		return
	}
	re = data
	return
}