package handlers

type HandlerResponse struct {
	replyMessages []string
}

func (h *HandlerResponse) AddReplyMessage(msg string) {
	h.replyMessages = append(h.replyMessages, msg)
}

func (h *HandlerResponse) GetReplyMessages() []string {
	result := make([]string, len(h.replyMessages))
	copy(result, h.replyMessages)
	return result
}
