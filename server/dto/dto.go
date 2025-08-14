package dto

type (
	BaseResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}

	VPNStatusResponse struct {
		Active bool   `json:"active"`
		Server string `json:"server"`
		Logs   string `json:"logs"`
	}

	VPNActivateRequest struct {
		Host   string `json:"host"`
		Port   string `json:"port"`
		Cookie string `json:"cookie"`
	}
)
