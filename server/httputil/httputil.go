package httputil

import (
	"encoding/json"
	"net/http"
	"server/dto"

	"github.com/rs/zerolog/log"
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
		log.Err(err).Msg("Write response failed")
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
		log.Err(err).Msg("Write response failed")
	}
}
