package PasteShare

type Paste struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Data        string `json:"data"`
	Date        string `json:"date"`
	Likes       int    `json:"likes"`
}

type UsersPastes struct {
	Id      int
	UserId  int
	PasteId int
}
