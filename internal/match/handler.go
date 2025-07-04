package match

type MatchHandler struct {
	service MatchService
}

func NewMatchHandler(service MatchService) *MatchHandler {
	return &MatchHandler{service: service}
}

func (h *MatchHandler) MatchEnded(result MatchResult) error {
	return h.service.HandleMatchResult(result)
}

func (h *MatchHandler) ItemUsed(item Item) error {
	return h.service.HandleUsedItem(item)
}
