package datamodels

type Product struct {
	Id           int64  `json:"id" sql:"id" form:"Id"`
	ProductName  string `json:"productName" sql:"productName" form:"productName"`
	ProductNum   int64  `json:"productNum" sql:"productNum" form:"productNum"`
	ProductImage string `json:"productImage" sql:"productImage" form:"productImage"`
	ProductUrl   string `json:"productUrl" sql:"productUrl" form:"productUrl"`
}
