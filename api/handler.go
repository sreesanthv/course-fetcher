package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/sreesanthv/course-fetcher/interfaces"
)

const SUCCESS_REPONSE int = 2

type Handler struct {
	logger *logrus.Logger
	store  interfaces.Store
}

func NewHandler(logger *logrus.Logger, store interfaces.Store) *Handler {
	return &Handler{
		logger: logger,
		store:  store,
	}
}

// read request
func (h *Handler) parseJSONBody(r *http.Request, reqData interface{}) error {
	buf := new(bytes.Buffer)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.logger.Error(err)
	}

	err = json.Unmarshal(buf.Bytes(), reqData)
	if err != nil {
		h.logger.Error(err)
	}

	return err
}

type responseData struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// send success repsone
func (h *Handler) sendResponse(w http.ResponseWriter, resData interface{}, message string) {
	dt := &responseData{
		Status:  "ok",
		Data:    resData,
		Message: message,
	}

	jData, err := json.Marshal(dt)
	if err != nil {
		h.logger.Error(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) badDataResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Invalid request data"
	}

	dt := &responseData{
		Status:  "nok",
		Message: message,
	}

	jData, err := json.Marshal(dt)
	if err != nil {
		h.logger.Error(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(jData)
}

func (h *Handler) ServerError(w http.ResponseWriter) {
	dt := &responseData{
		Status:  "nok",
		Message: "Failed to process the request",
	}

	jData, err := json.Marshal(dt)
	if err != nil {
		h.logger.Error(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jData)
}
