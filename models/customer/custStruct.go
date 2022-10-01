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

type Ordertb struct {
	Orderid      string  `form:"Orderid" json:"Customerid" `
	Visittime    string  `form:"Visittime" json:"Visittime" `
	Customerid   string  `form:"Customerid" json:"Customerid" `
	Consumetype  string  `form:"Consumetype" json:"Consumetype" `
	Item         string  `form:"Item" json:"Item" `
	Treatnum     int     `form:"Treatnum" json:"Treatnum" `
	Swmdmoney    float64 `form:"Swmdmoney" json:"Swmdmoney" `
	Qbmoney      float64 `form:"Qbmoney" json:"Qbmoney" `
	Mmmoney      float64 `form:"Mmmoney" json:"Mmmoney" `
	Mmboxnum     int     `form:"Mmboxnum" json:"Mmboxnum" `
	Donateboxnum int     `form:"Donateboxnum" json:"Donateboxnum" `
	Dgfmoney     float64 `form:"Dgfmoney" json:"Dgfmoney" `
	Xfymoney     float64 `form:"Xfymoney" json:"Xfymoney" `
	Totalmoney   float64 `form:"Totalmoney" json:"Totalmoney" `
	Returnmoney  float64 `form:"Returnmoney" json:"Returnmoney" `
	Kkmoney      float64 `form:"Kkmoney" json:"Kkmoney" `
	Owedmoney    float64 `form:"Owedmoney" json:"Owedmoney" `
	Paidmoney    float64 `form:"Paidmoney" json:"Paidmoney" `
	Consultteach string  `form:"Consultteach" json:"Consultteach" `
	Dqteach      string  `form:"Dqteach" json:"Dqteach" `
	Operateach   string  `form:"Operateach" json:"Operateach" `
	Operanum     int     `form:"Operanum" json:"Operanum" `
	Unoperanum   int     `form:"Unoperanum" json:"Unoperanum" `
	Comment      string  `form:"Comment" json:"Comment" `
}
