package datamodels

type Order struct {
	Id          int64 `sql:"id"`
	UserId      int64 `sql:"userId"`
	ProductId   int64 `sql:"productId"`
	OrderStatus int   `sql:"orderStatus"`
}

const (
	OrderWait = iota
	OrderSuccess
	OrderFailed
)
