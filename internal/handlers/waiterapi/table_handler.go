package waiterapi

import (
	"net/http"

	"CommonTable/internal/services/foodorderservices"
)

type TableHandler struct {
	tableService *foodorderservices.TableService
}

func NewTableHandler(tableService *foodorderservices.TableService) *TableHandler {
	return &TableHandler{tableService: tableService}
}

type tableJSON struct {
	ID             int32  `json:"id"`
	Number         int32  `json:"number"`
	Capacity       int32  `json:"capacity"`
	Status         string `json:"status"`
	HasActiveOrder bool   `json:"has_active_order"`
}

// List handles GET /api/waiter/tables
func (h *TableHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	tables, err := h.tableService.ListTablesWithActiveOrderFlag(r.Context(), claims.RestaurantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load tables")
		return
	}

	out := make([]tableJSON, 0, len(tables))
	for _, t := range tables {
		out = append(out, tableJSON{ID: t.ID, Number: t.Number, Capacity: t.Capacity, Status: string(t.Status), HasActiveOrder: t.HasActiveOrder})
	}

	writeJSON(w, http.StatusOK, out)
}
