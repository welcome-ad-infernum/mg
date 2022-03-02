package dto

type Target struct {
	ID      string      `json:"id"`
	URL     string      `json:"url"`
	Method  string      `json:"method"`
	Data    []byte      `json:"data"`
	Headers [][2]string `json:"headers"`
	Proxy   string      `json:"proxy_url"`
}

type TargetError struct {
	Target
	ErrCode int `json:"err_code"`
}
