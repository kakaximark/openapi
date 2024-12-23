package model

// Response 标准 API 响应
type Response struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Success"`
	Data    struct {
		Token string `json:"token" example:"eyJhbGciOiJ..."`
	} `json:"data"`
}
