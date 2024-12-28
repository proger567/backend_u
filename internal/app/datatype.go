package app

type User struct {
	Name       string `json:"user_name"`
	Role       string `json:"role_name"`
	RoleID     int    `json:"role_id"`
	CreateTime string `json:"create_time,omitempty"`
}

type Role struct {
	ID   int    `json:"id"`
	Role string `json:"role_name"`
}
