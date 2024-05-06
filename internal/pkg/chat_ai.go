package pkg

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Kdaito/accountant-line-bot/internal/types"
)

type ChatAI struct {
	apiUrl string
	apiKey string
}

func (c *ChatAI) ScanReceipt(file *os.File) {
	var messages []*types.ChatAIRequestMessage
	var contents []*types.ChatAiRequestContent

	base64Image, err := c.encodeImage(file)

	if err != nil {
		// TODO error処理
	}

	// プロンプトを追加
	contents = append(contents, &types.ChatAiRequestContent{
		Type: "text",
		Text: "What is in this image?",
	})

	// 画像を追加
	contents = append(contents, &types.ChatAiRequestContent{
		Type: "image_url",
		ImageURL: &types.ImageURL {
			URL: fmt.Sprintf("data:image/jpeg;base64,%s", base64Image),
		},
	})

	messages = append(messages, &types.ChatAIRequestMessage{
		Role:    "user",
		Content: contents,
	})

	requestBody := types.ChatAIRequest{
		Model:    "gpt-4-turbo",
		Messages: messages,
	}

	requestJSON, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("", c.apiUrl, bytes.NewBuffer(requestJSON))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + c.apiKey)

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()


	fmt.Print(response)
}

func (c *ChatAI) encodeImage(file *os.File) (string, error) {
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size)

	_, err := file.Read(buffer)

	if err != nil {
		return "", err
	}
	encodedString := base64.StdEncoding.EncodeToString(buffer)

	return encodedString, nil
}