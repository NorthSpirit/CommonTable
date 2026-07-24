package admin

import (
	"errors"
	"net/http"
	"strconv"

	"CommonTable/internal/repository"
	"CommonTable/internal/services/foodorderservices"
	"CommonTable/internal/templates/admin/adminshared"
	"CommonTable/internal/templates/admin/table"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type TableHandler struct {
	tableService *foodorderservices.TableService
}

func NewTableHandler(tableService *foodorderservices.TableService) *TableHandler {
	return &TableHandler{tableService: tableService}
}

// List handles GET /admin/tables
func (h *TableHandler) List(w http.ResponseWriter, r *http.Request) {
	tables, err := h.tableService.ListTablesForAdmin(r.Context())
	if err != nil {
		http.Error(w, "failed to load tables", http.StatusInternalServerError)
		return
	}

	views := make([]adminshared.TableView, 0, len(tables))
	for _, t := range tables {
		views = append(views, toTableView(t))
	}

	if err := table.TablesPage(views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Delete handles POST /admin/tables/{id}/delete
func (h *TableHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid table id", http.StatusBadRequest)
		return
	}
	if err := h.tableService.DeleteTable(r.Context(), int32(id)); err != nil {
		http.Error(w, "failed to delete table", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/tables", http.StatusFound)
}

// Restore handles POST /admin/tables/{id}/restore
func (h *TableHandler) Restore(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid table id", http.StatusBadRequest)
		return
	}
	if err := h.tableService.RestoreTable(r.Context(), int32(id)); err != nil {
		http.Error(w, "failed to restore table", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/tables", http.StatusFound)
}

// NewForm handles GET /admin/tables/new
func (h *TableHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	if err := table.TableFormPage(false, adminshared.TableView{}, "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Create handles POST /admin/tables/new
func (h *TableHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	number, errNum := strconv.Atoi(r.FormValue("number"))
	capacity, errCap := strconv.Atoi(r.FormValue("capacity"))
	if errNum != nil || errCap != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = table.TableFormPage(false, adminshared.TableView{}, "Number and capacity must be valid numbers.").Render(r.Context(), w)
		return
	}

	_, err := h.tableService.InsertTable(r.Context(), int32(number), int32(capacity), strPtrOrNil(r.FormValue("section")), repository.TableStatusFree)
	if err != nil {
		if isUniqueViolation(err) {
			w.WriteHeader(http.StatusBadRequest)
			_ = table.TableFormPage(false, adminshared.TableView{Number: int32(number), Capacity: int32(capacity)}, "Table number "+strconv.Itoa(number)+" already exists.").Render(r.Context(), w)
			return
		}
		http.Error(w, "failed to create table", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/tables", http.StatusFound)
}

// EditForm handles GET /admin/tables/{id}/edit
func (h *TableHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid table id", http.StatusBadRequest)
		return
	}

	dbTable, err := h.tableService.GetTableByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load table", http.StatusInternalServerError)
		return
	}

	if err := table.TableFormPage(true, toTableView(dbTable), "").Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Update handles POST /admin/tables/{id}/edit
func (h *TableHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid table id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	number, errNum := strconv.Atoi(r.FormValue("number"))
	capacity, errCap := strconv.Atoi(r.FormValue("capacity"))
	if errNum != nil || errCap != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = table.TableFormPage(true, adminshared.TableView{ID: int32(id)}, "Number and capacity must be valid numbers.").Render(r.Context(), w)
		return
	}

	if err := h.tableService.UpdateTable(r.Context(), int32(id), int32(number), int32(capacity), strPtrOrNil(r.FormValue("section"))); err != nil {
		if isUniqueViolation(err) {
			w.WriteHeader(http.StatusBadRequest)
			_ = table.TableFormPage(true, adminshared.TableView{ID: int32(id), Number: int32(number), Capacity: int32(capacity)}, "Table number "+strconv.Itoa(number)+" already exists.").Render(r.Context(), w)
			return
		}
		http.Error(w, "failed to update table", http.StatusInternalServerError)
		return
	}

	statusRaw := r.FormValue("status")
	if statusRaw != "" {
		if err := h.tableService.SetTableStatus(r.Context(), int32(id), repository.TableStatus(statusRaw)); err != nil {
			http.Error(w, "failed to update table status", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/admin/tables", http.StatusFound)
}
