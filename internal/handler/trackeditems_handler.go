package handler

import(
	"net/http"
	"github.com/al4nzh/pollingservice.git/internal/service"
)

type TrackedItemsHandler struct {
	service *service.TrackedItemsService
}

func NewTrackedItemsHandler(service *service.TrackedItemsService) *TrackedItemsHandler {
	return &TrackedItemsHandler{service: service}
}

func (h *TrackedItemsHandler) CreateTrackedItem(w http.ResponseWriter, r *http.Request){
	
}
