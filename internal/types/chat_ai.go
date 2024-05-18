package types

type ChatAIRequest struct {
	Model    string                  `json:"model"`
	Messages []*ChatAIRequestMessage `json:"messages"`
}

type ChatAIRequestMessage struct {
	Role    string                  `json:"role"`
	Content []*ChatAiRequestContent `json:"content"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type ChatAiRequestContent struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}
