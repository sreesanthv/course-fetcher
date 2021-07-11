package api

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/sreesanthv/course-fetcher/services"
)

type CourseHandler struct {
	Handler
	courseService *services.CourseService
}

func NewCourseHandler(handler *Handler) *CourseHandler {
	ah := &CourseHandler{
		Handler:       *handler,
		courseService: services.NewCourseService(handler.logger, handler.store),
	}

	return ah
}

func (h *CourseHandler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/fetch", h.fetch)
	r.Get("/search", h.search)
	return r
}

type fetchRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

func (h *CourseHandler) fetch(w http.ResponseWriter, r *http.Request) {
	req := &fetchRequest{}
	err := h.parseJSONBody(r, req)
	if err != nil {
		h.badDataResponse(w, "")
		return
	}

	go h.courseService.Fetch(req.Query, req.Limit)
	h.sendResponse(w, nil, "Course fetching has been started")
}

type searchReq struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type searchReply struct {
	List       []map[string]string `json:"list"`
	Pagination map[string]int      `json:"pagination"`
}

func (h *CourseHandler) search(w http.ResponseWriter, r *http.Request) {
	req := &searchReq{}
	err := h.parseJSONBody(r, req)
	if err != nil {
		h.badDataResponse(w, "")
		return
	} else if strings.TrimSpace(req.Query) == "" {
		h.badDataResponse(w, "Please input a query")
		return
	}

	if req.Limit == 0 {
		req.Limit = 10
	} else if req.Limit > 1000 {
		req.Limit = 1000
	}

	list, err := h.courseService.GetList(req.Query, req.Limit, req.Offset)
	if err != nil {
		h.ServerError(w)
		return
	}

	resp := searchReply{
		List: list,
		Pagination: map[string]int{
			"limit":  req.Limit,
			"offset": req.Offset,
			"count":  len(list),
		},
	}

	h.sendResponse(w, resp, "")
}
