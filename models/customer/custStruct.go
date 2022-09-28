package customer

import "time"

type Customertb struct {
	Customerid   string    `form:"customerid" json:"customerid" binding:"required"`
	Name         string    `form:"name" json:"name" binding:"required"`
	Gender       string    `form:"gender" json:"gender" binding:"required"`
	Consumetype  string    `form:"consumetype" json:"consumetype" binding:"required"`
	Visittime    time.Time `form:"visittime" json:"visittime" binding:"required"`
	Phone        string    `form:"phone" json:"phone" binding:"required"`
	Shop         string    `form:"shop" json:"shop" binding:"required"`
	Consultteach string    `form:"consultteach" json:"consultteach" binding:"required"`
	Item         string    `form:"item" json:"item" binding:"required"`
}
