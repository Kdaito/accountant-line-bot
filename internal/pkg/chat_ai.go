package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/Kdaito/accountant-line-bot/internal/lib/app_error"
	"github.com/Kdaito/accountant-line-bot/internal/types"
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
Please analyze the receipt image I will send next and output the information in the following JSON format. Extract the date, total amount, and list of items from the receipt. Follow the descriptions of each item to obtain the necessary information.

JSON format:
{
    "isReceipt": boolean,
    "date": "string",
    "totalAmount": number,
    "totalAmountIncludingTax": number,
    "currencySymbol": "string",
    "items": [
        {
            "name": "string",
            "amount": number,
            "count": number
        }
    ]
}

Description of each item:
- isReceipt: Whether the analyzed image is a receipt (true for a receipt, false for other images)
- date: The date on the receipt (format: YYYY-MM-DD)
- totalAmount: Total amount excluding tax (numeric only, without currency symbol)
- totalAmountIncludingTax: Total amount including tax (numeric only, without currency symbol)
- currencySymbol: Currency symbol
- items: Array of items
    - name: Name of the item
    - amount: Price of the item (numeric only, without currency symbol)
    - count: Quantity of the item

Notes:
- Omit items from the JSON if the information is not found.
- If no items are found, the "items" array should be empty.
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
type GptResponse struct {
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

func (c *ChatAI) ScanReceipt(encodedImage string) (*types.Receipt, error) {
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

	requestJSON, err := json.Marshal(requestBody)

	if err != nil {
		return nil, app_error.NewAppError(http.StatusInternalServerError, "Failed marshal request object for chat gpt api.", err)
	}

	req, err := http.NewRequest("POST", c.apiUrl, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, app_error.NewAppError(http.StatusInternalServerError, "Failed create request to chat gpt api.", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		return nil, app_error.NewAppError(http.StatusInternalServerError, "Failed request to chat gpt api.", err)
	}

	body, err := ioutil.ReadAll(response.Body)

	var res GptResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, app_error.NewAppError(http.StatusInternalServerError, "Failed unmarshal response of gpt api.", err)
	}

	defer response.Body.Close()

	// gptのレスポンスからjsonのみを取り出す
	re := regexp.MustCompile("```json\\n((?s:.*?))\\n```")
	match := re.FindStringSubmatch(res.Choices[0].Message.Content)
	jsonData := match[1]

	var result types.Receipt

	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		return nil, app_error.NewAppError(http.StatusInternalServerError, "Failed unmarshal json of recipt data.", err)
	}

	return &result, nil
}
