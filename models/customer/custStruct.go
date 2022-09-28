package customer

type Customertb struct {
	Customerid   string `form:"customerid" json:"customerid" binding:"required"`
	Name         string `form:"name" json:"name" binding:"required"`
	Gender       string `form:"gender" json:"gender" binding:"required"`
	Visittime    string `form:"visittime" json:"visittime" binding:"required"`
	Phone        string `form:"phone" json:"phone" binding:"required"`
	Shop         string `form:"shop" json:"shop" binding:"required"`
	Consultteach string `form:"consultteach" json:"consultteach" binding:"required"`
	Item         string `form:"item" json:"item" `
	Consumetype  string `form:"consumetype" json:"consumetype" `
}
