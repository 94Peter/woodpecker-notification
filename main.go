package main

import (
	"fmt"
	"os"

	"woodpecker-webhook/service"
)

func main() {
	sendMsgFuncSlice, err := service.GetSendMessageFunSlice()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	haserror := false
	for _, sendMsgFunc := range sendMsgFuncSlice {
		err = sendMsgFunc()
		if err != nil {
			fmt.Println(err)
			haserror = true
		}
	}
	if haserror {
		os.Exit(1)
	}
}
