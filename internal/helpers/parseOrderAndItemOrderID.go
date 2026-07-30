package helpers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func parseOrderAndItemOrderID(r *http.Request) (orderID int32, itemOrderID int32, ok bool) {
	oid, err1 := strconv.Atoi(chi.URLParam(r, "orderID"))
	iid, err2 := strconv.Atoi(chi.URLParam(r, "itemOrderID"))
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return int32(oid), int32(iid), true
}
