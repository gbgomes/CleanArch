package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gbgomes/GoExpert/CleanArch/internal/entity"
	"github.com/gbgomes/GoExpert/CleanArch/internal/usecase"
	"github.com/gbgomes/GoExpert/CleanArch/pkg/events"
)

type WebOrderHandler struct {
	EventDispatcher   events.EventDispatcherInterface
	OrderRepository   entity.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
}

type Error struct {
	Message string `json:"message"`
}

func NewWebOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreatedEvent events.EventInterface,
) *WebOrderHandler {
	return &WebOrderHandler{
		EventDispatcher:   EventDispatcher,
		OrderRepository:   OrderRepository,
		OrderCreatedEvent: OrderCreatedEvent,
	}
}

// Create Product godoc
// @Summary      Create order
// @Description  Create order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request     body      usecase.OrderInputDTO  true  "order request"
// @Success      201
// @Failure      500         {object}  Error
// @Router       /order [post]
func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ListAccounts godoc
// @Summary      List orders
// @Description  get all orders
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        page      query     string  false  "page number"
// @Param        rows      query     string  false  "number of rows per page"
// @Success      200       {array}   usecase.OrderListInputDTO
// @Failure      404       {object}  Error
// @Failure      500       {object}  Error
// @Router       /order [get]
func (h *WebOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	rowsParam := r.URL.Query().Get("rows")
	sort := r.URL.Query().Get("sort")
	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 0
	}
	rows, err := strconv.Atoi(rowsParam)
	if err != nil {
		rows = 0
	}
	dto := usecase.OrderListInputDTO{Sort: sort, Page: page, Offset: rows}

	listOrder := usecase.NewListOrdersUseCase(h.OrderRepository)
	output, err := listOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
