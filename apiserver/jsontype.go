package apiserver

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Data struct {
	User User `json:"user"`
}

type User struct {
	Userid  int64  `json:"userid"`
	Userame string `json:"username"`
	Email   string `json:"email"`
}

type LoginRequest struct {
	Userid         int64  `json:"userid"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	DeviceID       string `json:"device_id"`
	TurnstileToken string `json:"turnstileToken"`
}

type TurnstileVerify struct {
	Secret           string `json:"secret"`
	Response         string `json:"response"`
	Remoteip         string `json:"remoteip"`
	Idempontency_key string `json:"idempotency_key"`
}

type TurnstileResponse struct {
	Success     bool                   `json:"success"`
	Chllenge_ts string                 `json:"challenge_ts"`
	Hostname    string                 `json:"hostname"`
	ErrorCodes  []string               `json:"error-codes"`
	Action      string                 `json:"action"`
	Cdata       string                 `json:"cdata"`
	Metadata    map[string]interface{} `json:"metadata"`
}
