package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// ネイティブメッセージングのメッセージの先頭32bitはその後に続くUTF-8でエンコードされたJSONの長さ
	var length int32
	if err := binary.Read(os.Stdin, binary.LittleEndian, &length); err != nil {
		sendErrorToExtension(err)
		os.Exit(1)
	}

	message := make([]byte, length)
	if _, err := os.Stdin.Read(message); err != nil {
		sendErrorToExtension(err)
		os.Exit(1)
	}

	var dat map[string]any
	if err := json.Unmarshal(message, &dat); err != nil {
		sendErrorToExtension(err)
		os.Exit(1)
	}

	switch method := dat["method"].(string); method {
	case "openFolder", "openFile", "openFolderFromText", "openFileFromText":
		openPath(method, dat["path"].(string))
	case "showDiff":
		showDiff(dat["accessToken"].(string), dat["commandOptions"].(string), dat["url1"].(string), dat["versionNo1"].(string), dat["url2"].(string), dat["versionNo2"].(string))
	default:
		sendErrorStringToExtension("Unimplemented method.")
		os.Exit(1)
	}
}

func openPath(method string, varpath string) {
	path := replaceEnvVars(varpath)

	fmt.Fprintln(os.Stderr, path)

	info, err := os.Stat(path)
	if err != nil {
		sendErrorToExtension(err)
		os.Exit(1)
	}

	var cmd *exec.Cmd
	// 特定の環境(管理者権限もなく色々な制限があらかじめ加えられている環境)でexplorerを(CLIの)オプションなしでexecすると
	// フォルダーツリーを当該フォルダまで展開する(explorerのGUIの)オプションが有効になっているにも関わらず、
	// フォルダーツリーが展開されない事象が起きる。
	// 試行錯誤の結果、/rootオプションをつけるとなぜか解消されることがわかったのでワークアラウンドとして/rootをつける。
	// ちなみにexplorerの/rootオプションは昔のWindowsではchroot的なオプションだったはずだが、
	// Windows10やWindows11のexplorerで試してみるとそのような効果はないようだ。
	// なお /select オプションの効果(指定したファイルを選択した状態でフォルダを開く)をあきらめれば、
	// explorerではなくstartを用いることで同等の機能を実現することはできる。
	if method == "openFile" || info.IsDir() {
		// methodがopenFileの場合は、対象がファイルなら関連付けられたアプリで開き、
		// 対象がフォルダの場合はフォルダを開く。
		// methodがopenFolderで対象がフォルダの場合もそのフォルダを自体を開く。
		cmd = exec.Command("explorer.exe", "/root,", path)
	} else {
		// methodがopenFolderで対象がファイルの場合は、ファイルがあるフォルダを開く
		cmd = exec.Command("explorer.exe", "/select,/root,", path)
	}
	err = cmd.Run()
	if err != nil {
		sendErrorStringToExtension("Error in starting Explorer. However, this is usually completed successfully.")
		os.Exit(1)
	}

	// ここに到達することは多分ない
	sendSucessStringToExtension("Explorer started normally.")
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

func sendSucessStringToExtension(str string) {
	sendToExtension(map[string]any{
		"status":  "success",
		"message": str,
	})
}

func sendErrorStringToExtension(str string) {
	sendToExtension(map[string]any{
		"status":  "error",
		"message": str,
	})
}

func sendErrorToExtension(err error) {
	sendToExtension(map[string]any{
		"status":  "error",
		"message": err.Error(),
	})
}

func sendToExtension(response map[string]any) {
	resBytes, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	length := int32(len(resBytes))
	if err := binary.Write(os.Stdout, binary.LittleEndian, length); err != nil {
		panic(err)
	}

	if _, err := os.Stdout.Write(resBytes); err != nil {
		panic(err)
	}
}

func showDiff(accessToken, commandOptions, url1, versionNo1, url2, versionNo2 string) {
	filePath1, err := downloadFile(accessToken, url1, versionNo1)
	if err != nil {
		sendErrorToExtension(err)
		os.Exit(1)
	}

	filePath2, err := downloadFile(accessToken, url2, versionNo2)
	if err != nil {
		sendErrorToExtension(err)
		os.Exit(1)
	}

	// cmd.exeの二重引用符の処理が直観的でないのと
	// 将来的なWindows以外への対応も想定してひとまず簡易的に自前でパース
	args := splitString(commandOptions)
	args = append(args, filePath1)
	args = append(args, filePath2)
	cmd := exec.Command(args[0], args[1:]...)

	err = cmd.Start()
	if err != nil {
		sendErrorStringToExtension("Error in starting the diff tool.")
		os.Exit(1)
	}

	sendSucessStringToExtension("The diff tool started normally.")
}

func splitString(input string) []string {
	var result []string
	var currentToken string
	var insideQuotes bool

	for _, char := range input {
		switch char {
		case ' ':
			if !insideQuotes && currentToken != "" {
				result = append(result, currentToken)
				currentToken = ""
			} else if insideQuotes {
				currentToken += " "
			}
		case '"':
			insideQuotes = !insideQuotes
		default:
			currentToken += string(char)
		}
	}

	if currentToken != "" {
		result = append(result, currentToken)
	}

	return result
}

func downloadFile(accessToken, fileUrl, versionNo string) (string, error) {
	fileInfo, err := getFileDirInfo(fileUrl, accessToken)
	if err != nil {
		return "", err
	}

	versions, err := getVersions(accessToken, fileUrl)
	if err != nil {
		return "", err
	}

	versionList := makeVersionList(fileInfo, versions)

	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	parsedUrl, err := url.Parse(fileUrl)
	if err != nil {
		return "", err
	}

	host := parsedUrl.Hostname()

	cacheFilePath := filepath.Join(
		userCacheDir,
		"BoxUtilsHelper",
		"cache",
		host,
		"files",
		fileInfo["id"].(string),
		versionNo,
		fileInfo["name"].(string),
	)

	if fileExists(cacheFilePath) {
		// 既にダウンロード済みなのでダウンロード自体はスキップして正常応答する
		return cacheFilePath, nil
	}

	err = downloadContent(accessToken, fileUrl, versionList[versionNo], "", "", cacheFilePath)
	if err != nil {
		return "", err
	}

	return cacheFilePath, nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func getFileDirInfo(url, accessToken string) (map[string]any, error) {
	endpoint, err := getApiEndpoint(url)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getURLTypeAndIdentifier(url string) (string, string, error) {
	re := regexp.MustCompile(`https:\/\/.*\.?app.box.com\/([^/]+)\/(\d+)`)
	matches := re.FindStringSubmatch(url)
	if matches == nil || len(matches) != 3 {
		return "", "", errors.New("Invalid URL")
	}

	return matches[1], matches[2], nil
}

func getVersions(accessToken, url string) (map[string]any, error) {
	// https://ja.developer.box.com/reference/get-files-id-versions/
	// によると、
	// 「バージョンを追跡するのは、プレミアム (有償) アカウントを持つBoxユーザーのみです。」
	// とある。
	// しかし実際には無償アカウントでも取得できているように見える。
	// ただしバージョンのリストはoffsetを利用したとしても最新から100件までしか取得できないように見える。
	// またWebUIで表示されるV1やV2といったバージョン番号は
	// バージョンが100件を超えると最新の100件を指すようにスライドしていくように見えるため
	// バージョン番号とAPIで利用されるバージョンIDの対応は恒久的なものではない。
	// つまりキャッシュ不可であり、またWebUIとAPIの呼び出し時の間にファイル更新があった場合は対応が崩れる。
	endpoint, err := getVersionsApiEndpoint(url)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getApiEndpoint(taburl string) (string, error) {
	t, id, err := getURLTypeAndIdentifier(taburl)
	if err != nil || (t != "file" && t != "folder") {
		return "", errors.New("Invalid URL")
	}

	return fmt.Sprintf("https://api.box.com/2.0/%ss/%s", t, id), nil
}

func getContentApiEndpoint(taburl string) (string, error) {
	t, id, err := getURLTypeAndIdentifier(taburl)
	if err != nil || t != "file" {
		return "", errors.New("Invalid URL")
	}

	return fmt.Sprintf("https://api.box.com/2.0/%ss/%s/content", t, id), nil
}

func getVersionsApiEndpoint(taburl string) (string, error) {
	t, id, err := getURLTypeAndIdentifier(taburl)
	if err != nil || t != "file" {
		return "", errors.New("Invalid URL")
	}

	return fmt.Sprintf("https://api.box.com/2.0/%ss/%s/versions", t, id), nil
}

func downloadContent(accessToken, url, versionId, start, length, filepath string) error {
	// 対象が空ファイルの場合、Rangeヘッダを設定するとステータスコード416が返ってくるので注意

	endpoint, err := getContentApiEndpoint(url)
	if err != nil {
		return err
	}
	endpoint += "?version=" + versionId

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = saveToFile(res.Body, filepath)
	if err != nil {
		return err
	}

	return nil
}

func saveToFile(data io.Reader, filePath string) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	if err != nil {
		return err
	}

	return nil
}

func makeVersionList(file map[string]any, versions map[string]any) map[string]string {
	result := make(map[string]string)

	totalCount := int(versions["total_count"].(float64))

	if totalCount > 0 {
		entries := versions["entries"].([]any)
		for index := range entries {
			key := fmt.Sprintf("V%d", totalCount-index)
			entry := entries[index].(map[string]any)
			result[key] = entry["id"].(string)
		}
	}

	lastKey := fmt.Sprintf("V%d", totalCount+1)
	result[lastKey] = file["file_version"].(map[string]any)["id"].(string)

	return result
}
