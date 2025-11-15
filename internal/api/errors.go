package api

import (
	"encoding/json"
	"net/http"

	"github.com/guverz/pr-reviewer-service/internal/domain"
)

func WriteError(w http.ResponseWriter, err error) {
	domainErr, ok := err.(domain.DomainError)
	if !ok {
		// Неизвестная ошибка
		writeErrorResponse(w, http.StatusInternalServerError, domain.ErrorCodeNotFound, err.Error())
		return
	}

	// Маппинг кодов ошибок на HTTP статусы
	statusCode := http.StatusInternalServerError
	switch domainErr.Code {
	case domain.ErrorCodeTeamExists, domain.ErrorCodePRExists:
		statusCode = http.StatusBadRequest
	case domain.ErrorCodeNotFound:
		statusCode = http.StatusNotFound
	case domain.ErrorCodePRMerged, domain.ErrorCodeNotAssigned, domain.ErrorCodeNoCandidate:
		statusCode = http.StatusConflict
	default:
		statusCode = http.StatusInternalServerError
	}

	writeErrorResponse(w, statusCode, domainErr.Code, domainErr.Message)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, code domain.ErrorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    string(code),
			Message: message,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}



