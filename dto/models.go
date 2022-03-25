package dto

type Target struct {
	ID      int         `json:"id"`
	URL     string      `json:"url"`
	Method  string      `json:"method"`
	Data    []byte      `json:"data"`
	Headers [][2]string `json:"headers"`
	Proxy   string      `json:"proxy_url"`
}

type TargetResponse struct {
	Target
	Code int `json:"http_code"`
}

type Statistic struct {
	Success int64 `json:"success"`
	Error   int64 `json:"error"`
}

type TargetStatistic struct {
	Statistic
	AgentUID string `json:"agent"`
	TargetID int    `json:"-"`
}
