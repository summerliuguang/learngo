package apiserver

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
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
	Userid   int64  `json:"userid"`
	Username string `json:"username"`
	Password string `json:"password"`
	DeviceID string `json:"device_id"`
}
