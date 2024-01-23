package PasteShare

type Paste struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Data        string `json:"data" db:"data" binding:"required"`
	Date        string `json:"paste_date" db:"paste_date"`
	Likes       int    `json:"likes" db:"likes"`
}

type UsersPastes struct {
	Id      int
	UserId  int
	PasteId int
}
