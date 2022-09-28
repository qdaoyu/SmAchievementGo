package customer

type Customertb struct {
	Customerid   string `form:"customerid" json:"Customerid" binding:"required"`
	Name         string `form:"name" json:"Name" binding:"required"`
	Gender       string `form:"gender" json:"Gender" binding:"required"`
	Visittime    string `form:"visittime" json:"Visittime" binding:"required"`
	Phone        string `form:"phone" json:"Phone" binding:"required"`
	Shop         string `form:"shop" json:"Shop" binding:"required"`
	Consultteach string `form:"consultteach" json:"Consultteach" binding:"required"`
	// Item         string `form:"item" json:"item" `
	// Consumetype  string `form:"consumetype" json:"consumetype" `
}
