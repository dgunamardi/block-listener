package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"earhart.com/config"
)

// ToJson
var (
	vars config.ConfigVars
	cwd  string

	oldDir = "/home/tkgoh/Sandbox/block-listener"
	newDir string
)

const (
	configPath = "/ccp/bcs-test-channel-sdk-config.yaml"
	orgId      = "4f08db41ded98093a7266580a4a2ae3ce62ce74a"
	userName   = "Admin"
)

func main() {
	cwd, _ = os.Getwd()
	newDir = cwd

	SetConfigPath(configPath)
	SetOrgAndUser(orgId, userName)
	WriteToJson()

	fmt.Println(cwd)
	SetChannelConfigFile()
}

func SetConfigPath(relativePath string) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	vars.ChannelConfigPath = cwd + relativePath
}

func SetOrgAndUser(orgId string, userId string) {
	vars.OrgId = orgId
	vars.UserName = userId
}

func WriteToJson() {
	varsToJson, _ := json.Marshal(vars)

	err := ioutil.WriteFile(cwd+"/config/vars.json", varsToJson, 0644)
	if err != nil {
		panic(fmt.Errorf("f to write json vars: %v", err))
	}

}

func SetChannelConfigFile() {
	configIn, err := os.Open(vars.ChannelConfigPath)
	if err != nil {
		panic(err)
	}
	defer configIn.Close()

	configOut, err := os.OpenFile(vars.ChannelConfigPath, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		panic(err)
	}
	defer configOut.Close()

	//test, _ := ioutil.ReadFile(vars.ChannelConfigPath)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(test))

	br := bufio.NewReader(configIn)
	index := 1
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read err:", err)
			os.Exit(-1)
		}
		newLine := strings.Replace(string(line), oldDir, newDir, -1)
		//fmt.Println(newLine)
		_, err = configOut.WriteString(newLine + "\n")
		if err != nil {
			fmt.Println("write to file fail:", err)
			os.Exit(-1)
		}
		//fmt.Println("done ", index)
		index++
	}
	fmt.Println("FINISH!")
}
