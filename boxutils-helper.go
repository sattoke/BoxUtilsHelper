package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	var length int32
	if err := binary.Read(os.Stdin, binary.LittleEndian, &length); err != nil {
		panic(err)
	}

	message := make([]byte, length)
	if _, err := os.Stdin.Read(message); err != nil {
		panic(err)
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(message, &dat); err != nil {
		panic(err)
	}

	response := map[string]interface{}{
		"response": fmt.Sprintf("Hello, %s", dat["path"]),
	}

	resBytes, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	length = int32(len(resBytes))
	if err := binary.Write(os.Stdout, binary.LittleEndian, length); err != nil {
		panic(err)
	}

	if _, err := os.Stdout.Write(resBytes); err != nil {
		panic(err)
	}

	path := replaceEnvVars(dat["path"].(string))

	fmt.Fprintln(os.Stderr, path)

	info, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var cmd *exec.Cmd
	if dat["method"] == "openFile" || info.IsDir() {
		// methodがopenFileの場合は、対象がファイルなら関連付けられたアプリで開き、
		// 対象がフォルダの場合はフォルダを開く。
		// methodがopenFolderで対象がフォルダの場合もそのフォルダを自体を開く。
		cmd = exec.Command("explorer.exe", path)
	} else {
		// methodがopenFolderで対象がファイルの場合は、ファイルがあるフォルダを開く
		cmd = exec.Command("explorer.exe", "/select,", path)
	}
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		// explorerではなぜか必ずエラーになるので正常終了にしてしまう
		os.Exit(0)
	}

	fmt.Println("ok")
	os.Exit(0)
}

func replaceEnvVars(str string) string {
	pattern := regexp.MustCompile(`%([^%]+)%`)

	return pattern.ReplaceAllStringFunc(str, func(s string) string {
		envVar := strings.TrimPrefix(s, "%")
		envVar = strings.TrimSuffix(envVar, "%")

		envValue := os.Getenv(envVar)

		if envValue != "" {
			return envValue
		}

		// 環境変数が存在しない場合は置換しない
		return s
	})
}
