package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// PostLineMessage LINEにメッセージを送る
func PostLineMessage(postText string) error {
	client := &http.Client{}
	jsonStr := `{"messages":[{"type":"text","text":"` + postText + `"}]}`
	fmt.Println(jsonStr)
	req, err := http.NewRequest(
		"POST",
		"https://api.line.me/v2/bot/message/broadcast",
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "#####################################")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()

	return nil
}

// CreatePostMessage LINEに送信する文言を作成
func CreatePostMessage(serialCodeList []string, targetURL string) string {
	postMessage := "新しいシリアルコードが追加されました！\\n"

	for _, serialCode := range serialCodeList {
		postMessage += "【" + serialCode + "】\\n"
	}
	postMessage += targetURL

	return postMessage
}
