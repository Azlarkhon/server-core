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

func (sp *MatchHandler) ItemUsed(event InventoryEvent) error {
	return sp.service.HandleInventoryEvent(event)
}
