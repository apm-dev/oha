package httputils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/apm-dev/oha/src/domain"
	log "github.com/sirupsen/logrus"
)

type httpJsonErr struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Params  interface{} `json:"params,omitempty"`
}

func RespondWithErrJson(
	w http.ResponseWriter,
	err error,
	methodName string,
	errMsgReplacement string,
	params interface{},
) {
	status, msg := errToHttpStatusCodeAndMessage(err, methodName)
	resp := &httpJsonErr{
		Code:    status,
		Message: msg,
		Params:  params,
	}
	if errMsgReplacement != "" {
		resp.Message = errMsgReplacement
	}
	b, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("failed to marshal HttpJsonErr, err: '%s'", err)
		respondWithDataJSON(w, http.StatusInternalServerError, []byte(""))
		return
	}
	respondWithDataJSON(w, status, b)
}

func RespondWithStructJSON(
	w http.ResponseWriter,
	httpStatus int,
	payload interface{},
) {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("failed to marshal data, err: %s", err)
	}
	respondWithDataJSON(w, httpStatus, body)
}

func respondWithDataJSON(
	w http.ResponseWriter,
	httpStatus int,
	payload []byte,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_, err := w.Write(payload)
	if err != nil {
		log.Debugf("Failed to send out reply, error: %s", err)
	}
}

func errToHttpStatusCodeAndMessage(err error, methodName string) (int, string) {
	switch {
	case errors.Is(err, domain.ErrInvalidArgument):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, err.Error()
	default:
		return http.StatusInternalServerError, fmt.Sprintf("failed to %s", methodName)
	}
}
