package model

type User struct {
	Id          int
	UserName    string     `json:"user_name"`
	PasswordStr string     `json:"password"`
	Password    string     `json:"-"`
	Name        string     `json:"name"`
	Age         int        `json:"age"`
	Gender      string     `json:"gender"`
	PhoneNumber string     `json:"phone_number"`
	EmailID     string     `json:"email_id"`
	CreatedTime string `json:"created_time"`
	UpdatedTime string `json:"updated_time"`
	Token       string     `json:"-"`
}
