package httputil

import (
	"encoding/json"
	"log"
	"net/http"
	"server/dto"
)

func ResponseOK(w http.ResponseWriter, data any) {
	w.WriteHeader(http.StatusOK)

	response := dto.BaseResponse{
		Success: true,
		Message: "Success",
		Data:    data,
	}

	r, err := json.Marshal(response)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, err = w.Write(r)
	if err != nil {
		log.Default().Printf("Write response failed: %v", err)
	}
}

func ResponseError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	response := dto.BaseResponse{
		Success: false,
		Message: err.Error(),
	}

	r, err := json.Marshal(response)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	_, err = w.Write(r)
	if err != nil {
		log.Default().Printf("Write response failed: %v", err)
	}
}
