package repositories

import (
	"database/sql"
	"mall-seckill/common"
	"mall-seckill/datamodels"
	"strconv"
)

type IProduct interface {
	Conn() error
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: db}
}

func (p *ProductManager) Conn() error {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return nil
}

func (p *ProductManager) Insert(product *datamodels.Product) (int64, error) {
	if err := p.Conn(); err != nil {
		return 0, err
	}
	sql := "insert into " + p.table + " set productName=?,productNum=?,productImage=?,productUrl=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (p *ProductManager) Delete(productId int64) bool {
	if err := p.Conn(); err != nil {
		return false
	}
	sql := "delete from " + p.table + " where id=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return false
	}
	_, err = stmt.Exec(productId)
	if err != nil {
		return false
	}
	return true
}

func (p *ProductManager) Update(product *datamodels.Product) error {
	if err := p.Conn(); err != nil {
		return err
	}
	sql := "update " + p.table + " set productName=?,productNum=?,productImage=?,productUrl=? where id=" + strconv.FormatInt(product.Id, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	return err
}

func (p *ProductManager) SelectByKey(productID int64) (*datamodels.Product, error) {
	if err := p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}
	sql := "Select * from " + p.table + " where id=" + strconv.FormatInt(productID, 10)
	row, err := p.mysqlConn.Query(sql)
	defer row.Close()
	if err != nil {
		return &datamodels.Product{}, err
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}
	product := &datamodels.Product{}
	common.DataToStructByTagSql(result, product)
	return product, nil
}

func (p *ProductManager) SelectAll() ([]*datamodels.Product, error) {
	if err := p.Conn(); err != nil {
		return []*datamodels.Product{}, err
	}
	sql := "Select * from " + p.table
	rows, err := p.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return []*datamodels.Product{}, err
	}
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return []*datamodels.Product{}, nil
	}
	var productArray []*datamodels.Product
	for _, v := range result {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}
	return productArray, nil
}
