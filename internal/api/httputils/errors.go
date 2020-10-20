package httputils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/musicmash/artisync/internal/log"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteInternalError(w http.ResponseWriter, err error) {
	log.Error(err.Error())
	WriteErrorWithCode(w, http.StatusInternalServerError, errors.New("internal")) //nolint:goerr113
}

func WriteClientError(w http.ResponseWriter, err error) {
	WriteErrorWithCode(w, http.StatusBadRequest, err)
}

func WriteErrorWithCode(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	b, _ := json.Marshal(&ErrorResponse{Error: err.Error()})
	_, _ = w.Write(b)
}
