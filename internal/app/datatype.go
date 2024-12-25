package app

type User struct {
	Name       string `json:"user_name"`
	Role       string `json:"role_name"`
	CreateTime string `json:"create_time,omitempty"`
}
