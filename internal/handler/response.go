package handler

type ListSentMessagesResponse struct {
	Messages []MessageResponse `json:"messages"`
	Total    int               `json:"total"`
}

type MessageResponse struct {
	ID        int    `json:"id"`
	To        string `json:"to"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	SentAt    string `json:"sent_at,omitempty"`
	MessageID string `json:"message_id,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ToggleSchedulerResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
