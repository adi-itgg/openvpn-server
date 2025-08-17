package controller

import (
	"fmt"
	"io"
	"net/http"
	"server/dto"
	"server/httputil"
	"server/usecase"

	"github.com/goccy/go-json"
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

func (c *Controller) Servers(w http.ResponseWriter, r *http.Request) {
	data, err := c.usecase.Servers()
	if err != nil {
		httputil.ResponseError(w, err)
		return
	}
	httputil.ResponseOK(w, data)
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
		fmt.Printf("Read body failed: %v\n", err)
		httputil.ResponseError(w, err)
		return
	}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		fmt.Printf("Marshal body failed: %v\n", err)
		httputil.ResponseError(w, err)
		return
	}

	err = c.usecase.Activate(&body)
	if err != nil {
		httputil.ResponseError(w, err)
		return
	}
}
