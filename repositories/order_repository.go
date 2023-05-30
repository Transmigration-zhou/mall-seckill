package repositories

import (
	"database/sql"
	"mall-seckill/common"
	"mall-seckill/datamodels"
	"strconv"
)

type IOrder interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

type OrderManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewOrderManager(table string, sql *sql.DB) IOrder {
	return &OrderManager{table: table, mysqlConn: sql}
}

func (o *OrderManager) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

func (o *OrderManager) Insert(order *datamodels.Order) (int64, error) {
	if err := o.Conn(); err != nil {
		return 0, err
	}
	sql := "insert into " + o.table + " set userId=?,productID=?,orderStatus =?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (o *OrderManager) Delete(orderId int64) bool {
	if err := o.Conn(); err != nil {
		return false
	}
	sql := "delete from " + o.table + " where id=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return false
	}
	_, err = stmt.Exec(orderId)
	if err != nil {
		return false
	}
	return true
}

func (o *OrderManager) Update(order *datamodels.Order) error {
	if err := o.Conn(); err != nil {
		return err
	}
	sql := "update " + o.table + " set userId=?,productId=?,orderStatus=? where id=" + strconv.FormatInt(order.Id, 10)
	stmt, err := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	return err
}

func (o *OrderManager) SelectByKey(orderID int64) (*datamodels.Order, error) {
	if errConn := o.Conn(); errConn != nil {
		return &datamodels.Order{}, errConn
	}
	sql := "select * from " + o.table + " where id=" + strconv.FormatInt(orderID, 10)
	row, err := o.mysqlConn.Query(sql)
	defer row.Close()
	if err != nil {
		return &datamodels.Order{}, err
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, nil
	}
	order := &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return order, nil
}

func (o *OrderManager) SelectAll() ([]*datamodels.Order, error) {
	if err := o.Conn(); err != nil {
		return []*datamodels.Order{}, err
	}
	sql := "select * from " + o.table
	rows, err := o.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return []*datamodels.Order{}, err
	}
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return []*datamodels.Order{}, nil
	}
	var orderArray []*datamodels.Order
	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return orderArray, nil
}

func (o *OrderManager) SelectAllWithInfo() (map[int]map[string]string, error) {
	if err := o.Conn(); err != nil {
		return map[int]map[string]string{}, err
	}
	sql := "select o.id,p.productName,o.orderStatus from `order` as o left join product as p on o.productId=p.id"
	rows, err := o.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return map[int]map[string]string{}, err
	}
	return common.GetResultRows(rows), err
}
