package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
// serviceをメンバに持つTODOHandlerを新しく作って返しているんだね
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		req := &model.CreateTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if req.Subject == "" {
			http.Error(w, "Subject is empty", http.StatusBadRequest)
			return
		}
		res, err := h.Create(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}

	}

	if r.Method == http.MethodGet {
		req := &model.ReadTODORequest{}

		q := r.URL.Query()
		prevID := q.Get("prev_id")
		size := q.Get("size")

		var parsedPrevID int64
		var parsedSize int64

		log.Printf("prev_id: %s, size: %s", prevID, size)

		// if prevID == "" {
		// 	parsedPrevID = -1
		// } else {
		// 	parsedPrevID, _ = strconv.ParseInt(prevID, 10, 64)
		// }
		parsedPrevID, _ = strconv.ParseInt(prevID, 10, 64)

		parsedSize, _ = strconv.ParseInt(size, 10, 64)

		log.Print("parsedPrevID: ", parsedPrevID)

		req.PrevID, req.Size = parsedPrevID, parsedSize

		res, err := h.Read(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPut {
		req := &model.UpdateTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if req.ID == 0 || req.Subject == "" {
			http.Error(w, "Bad Request!", http.StatusBadRequest)
			return
		}
		res, err := h.Update(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	// if r.Method == http.MethodDelete {
	// }
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{
		TODO: *todo,
	}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todo, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{
		TODOs: todo,
	}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{
		TODO: *todo,
	}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
