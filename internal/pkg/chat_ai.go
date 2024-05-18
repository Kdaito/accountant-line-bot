package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/Kdaito/accountant-line-bot/internal/lib/app_error"
	// "github.com/Kdaito/accountant-line-bot/internal/types"
)

type ChatAI struct {
	apiUrl string
	apiKey string
}

func NewChatAI(apiUrl, apiKey string) *ChatAI {
	return &ChatAI{apiUrl: apiUrl, apiKey: apiKey}
}

// Chat GPT APIに送るプロンプト
const prompt = `
Please analyze the receipt image that will be sent next and return it in JSON format according to the following structure. Extract the date, total amount, and list of items from the receipt.

JSON format:
{
		"date": "string",
		"totalAmount": "string",
		"items": [
				{
						"name": "string",
						"amount": "string",
						"count": "string"
				}
		]
}

Description of each field:
- date: The date on the receipt (format: YYYY-MM-DD)
- totalAmount: The total amount (including currency symbol)
- items: Array of items
		- name: Name of the item
		- amount: Amount of the item (including currency symbol)
		- count: Quantity of the item

Notes:
- If the date is not found, set the "date" field to "N/A".
- If the total amount is not found, set the "totalAmount" field to "N/A".
- If no items are found, set the "items" array to empty.
`

type ChatMessage struct {
	Role    string                 `json:"role"`
	Content []ChatAiRequestContent `json:"content"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type ChatAiRequestContent struct {
	Type     string   `json:"type"`
	Text     string   `json:"text,omitempty"`
	ImageURL ImageURL `json:"image_url,omitempty"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// チャットの受信データ
type Response struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Model   string    `json:"model"`
	Usage   *Usage    `json:"usage"`
	Choices []*Choice `json:"choices"`
}

// APIの使用量
type Usage struct {
	// 入力データのトークン
	PromptTokens int `json:"prompt_tokens"`

	// 出力データのトークン
	CompletionTokens int `json:"completion_tokens"`

	// 合計トークン
	TotalTokens int `json:"total_tokens"`
}

type Choice struct {
	// 受信メッセージ
	Message *ResponseMessage `json:"message"`

	// リクエストが異常終了した場合の理由(正常終了の場合は空文字)
	FinishReason string `json:"finish_reason"`

	// トークン化されたインデックス
	Index int `json:"index"`
}

// チャットの受信メッセージ
type ResponseMessage struct {
	// メッセージの役割(assistant, user, systemのどれか)
	Role string `json:"role"`

	// メッセージの本文
	Content string `json:"content"`
}

func (c *ChatAI) ScanReceipt(encodedImage string) error {
	var contents []ChatAiRequestContent

	contents = append(contents, ChatAiRequestContent{
		Type: "text",
		Text: prompt,
	})

	contents = append(contents, ChatAiRequestContent{
		Type: "image_url",
		ImageURL: ImageURL{
			URL: fmt.Sprintf("data:image/jpeg;base64,%s", encodedImage),
		},
	})

	message := ChatMessage{
		Role:    "user",
		Content: contents,
	}

	requestBody := ChatRequest{
		Model:    "gpt-4o", // 必要に応じてモデルを変更してください
		Messages: []ChatMessage{message},
	}

	requestJSON, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", c.apiUrl, bytes.NewBuffer(requestJSON))
	if err != nil {
		return app_error.NewAppError(http.StatusInternalServerError, "Failed create request to chat gpt api.", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{}

	fmt.Print("gptにリクエストを送信します。")
	response, err := client.Do(req)

	if err != nil {
		return app_error.NewAppError(http.StatusInternalServerError, "Failed request to chat gpt api.", err)
	}

	body, err := ioutil.ReadAll(response.Body)

	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil
	}

	defer response.Body.Close()

	re := regexp.MustCompile("```json\\n((?s:.*?))\\n```")
	match := re.FindStringSubmatch(res.Choices[0].Message.Content)
	jsonData := match[1]
	fmt.Println(jsonData)

	// 以下jsonData(例)
	// 	{
	//     "date": "2024-04-12",
	//     "totalAmount": "€9,36",
	//     "items": [
	//         {
	//             "name": "CAMP YOGHURT",
	//             "amount": "€3,19",
	//             "count": "1"
	//         },
	//         {
	//             "name": "DOOSJE FRUIT",
	//             "amount": "€2,59",
	//             "count": "1"
	//         },
	//         {
	//             "name": "RIBBELCHIPS",
	//             "amount": "€1,39",
	//             "count": "1"
	//         },
	//         {
	//             "name": "SCHARRELEI",
	//             "amount": "€2,19",
	//             "count": "1"
	//         }
	//     ]
	// 	}
	return nil
}
