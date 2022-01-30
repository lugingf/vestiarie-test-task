package handler

import (
	"encoding/json"
	"fmt"
	"github.com/lugingf/vestiarie-test-task/internal/domain"
	"github.com/lugingf/vestiarie-test-task/resources"
	"log"
	"net/http"
)

type PayoutHandler struct {
	Service *domain.PayoutService
}

func NewPayoutHandler(s *domain.PayoutService) (*PayoutHandler, error) {
	return &PayoutHandler{Service: s}, nil
}

func (h *PayoutHandler) PostPayouts(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	updateId := r.Header.Get("PaymentUpdateID")
	if updateId == "" {
		log.Printf("No header PaymentUpdateID")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No header PaymentUpdateID"))
	}

	items := make([]domain.Item, 0)
	err := decoder.Decode(&items)

	if err != nil {
		log.Printf("Cant get body. Error: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cant get request body"))
	}

	payouts, err := h.Service.StorePayouts(items, updateId)
	if err == resources.ErrUpdateIdExists {
		log.Printf("PaymentUpdateID: %s Error: %v",updateId,  err.Error())
		w.Write([]byte(fmt.Sprintf("PaymentUpdateID: %s is already exists",updateId)))
		return
	}
	if err != nil {
		log.Printf("Cant save payouts. Error: %v", err.Error())
		h.writeCommonError(w)
		return
	}

	payoutsJson, err := json.Marshal(payouts)
	if err != nil {
		h.writeCommonError(w)
		return
	}

	w.Write(payoutsJson)
}

func (h *PayoutHandler) writeCommonError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}