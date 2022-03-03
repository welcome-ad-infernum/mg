package dto

type Target struct {
	ID      int         `json:"id"`
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

type Statistic struct {
	AgentUID string `json:"agent"`
	Success  int64  `json:"success"`
	Error    int64  `json:"error"`
}

type TargetStatistic struct {
	Statistic
	TargetID int `json:"target_id"`
}
