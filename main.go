package WSATRAINING
// package main

import (
	"fmt"
	"log"
	"os"

	"wsatraining/api"
	"wsatraining/config"
	"wsatraining/utils"
)

func main() {
	utils.Logger("")
	config.Init(nil) //Init Config.yaml
	api.StartRouter()
}

//Start ...
func Start(OS, OSPATH, PROXY, ADDRESS string) {
	//OS e.g "ios" or "android"
	//PATH e.g "/sdcard/com.sample.app/"
	var yaml = []byte(fmt.Sprintf(`os: %v
path: %v
proxy: %v
address: %v`, OS, OSPATH, PROXY, ADDRESS))

	utils.Logger(OSPATH)
	config.Init(yaml) //Init Config.yaml
	go api.StartRouter()
}

//Stop ...
func Stop() {
	sMessage := "stopping service @ " + config.Get().Address
	println(sMessage)
	log.Println(sMessage)
	os.Exit(1)
}
