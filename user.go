package PasteShare

type User struct {
	Id       int    `json:"id" db:"id"`
	NickName string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
