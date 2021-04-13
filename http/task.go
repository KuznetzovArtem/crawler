package http

import (
	"crawler/domain"
	"crawler/repositories/logger"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"runtime/debug"
)

// taskHandler service for handling transport for tasks.
type taskHandler struct {
	encoder *encoder
	log     logger.Logger
	tasks   TaskSender
}

type TaskSender interface {
	Do(r domain.Request) (*domain.Response, error)
}

func newTaskHandler(encoder *encoder, tasks TaskSender) *taskHandler {
	return &taskHandler{
		encoder: encoder,
		tasks:   tasks,
	}
}

func (h *taskHandler) Routes(router chi.Router) {
	router.Get("/", h.createTask)
}

// createTask this is method which preparing income and out data
func (h *taskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	defer func(request *http.Request) {
		if msg := recover(); msg != nil {
			if len(requestBody) > 0 {
				go h.log.Log(logger.Panic, fmt.Sprintf("panic in createTask func: %v %v %v", msg, string(debug.Stack()), string(requestBody)))
			} else {
				go h.log.Log(logger.Panic, fmt.Sprintf("empty body request, panic in createTask func: %v %v", msg, string(debug.Stack())))
			}

		}
	}(r)
	jsonResponse, err := h.sendRequest(requestBody)
	if err != nil {
		go h.log.Log(logger.Error, err.Error())
		h.encoder.Error(w, err, http.StatusInternalServerError)
		return
	}
	h.encoder.StatusResponse(w, jsonResponse, http.StatusOK)
	return
}

// sendRequest - function send request to special task
func (h *taskHandler) sendRequest(requestBody []byte) ([]byte, error) {
	var (
		response *domain.Response
		err      error
		request  domain.Request
	)
	err = json.Unmarshal(requestBody, &request)
	if err != nil {
		return nil, err
	}
	if response, err = h.tasks.Do(request); err != nil {
		return nil, err
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return jsonResponse, nil
}
