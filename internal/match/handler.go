package match

type MatchHandler struct {
	service MatchService
}

func NewMatchHandler(service MatchService) *MatchHandler {
	return &MatchHandler{service: service}
}

func (sp *MatchHandler) MatchEnded(result MatchResult) error {
	return sp.service.HandleMatchResult(result)
}

func (sp *MatchHandler) ItemUsed(item Item) error {
	return sp.service.HandleUsedItem(item)
}
