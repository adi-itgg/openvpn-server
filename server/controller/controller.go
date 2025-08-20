package controller

import (
	"io"
	"net/http"
	"server/dto"
	"server/httputil"
	"server/usecase"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
)

func NewController(
	usecase *usecase.Usecase,
) *Controller {
	return &Controller{
		usecase: usecase,
	}
}

type Controller struct {
	usecase *usecase.Usecase
}

func (c *Controller) Status(w http.ResponseWriter, r *http.Request) {
	data, err := c.usecase.Status()
	if err != nil {
		httputil.ResponseError(w, err)
		return
	}
	httputil.ResponseOK(w, data)
}

func (c *Controller) Activate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var body dto.VPNActivateRequest

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Err(err).Msg("Read body failed")
		httputil.ResponseError(w, err)
		return
	}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		log.Err(err).Msg("Marshal body failed")
		httputil.ResponseError(w, err)
		return
	}

	err = c.usecase.Activate(&body)
	if err != nil {
		httputil.ResponseError(w, err)
		return
	}
	httputil.ResponseOK(w, nil)
}
