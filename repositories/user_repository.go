package repositories

import (
	"database/sql"
	"errors"
	"mall-seckill/common"
	"mall-seckill/datamodels"
	"strconv"
)

type IUser interface {
	Conn() error
	Select(string) (*datamodels.User, error)
	Insert(*datamodels.User) (int64, error)
	SelectByID(int64) (*datamodels.User, error)
}

type UserManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewUserManager(table string, mysqlConn *sql.DB) IUser {
	return &UserManager{table, mysqlConn}
}

func (u *UserManager) Conn() error {
	if u.mysqlConn == nil {
		db, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlConn = db
	}
	if u.table == "" {
		u.table = "user"
	}
	return nil
}

func (u *UserManager) Select(userName string) (*datamodels.User, error) {
	if userName == "" {
		return &datamodels.User{}, errors.New("用户名不能为空")
	}
	if err := u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "select * from " + u.table + " where userName=?"
	rows, err := u.mysqlConn.Query(sql, userName)
	defer rows.Close()
	if err != nil {
		return &datamodels.User{}, err
	}
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在")
	}
	user := &datamodels.User{}
	common.DataToStructByTagSql(result[0], user)
	return user, nil
}

func (u *UserManager) Insert(user *datamodels.User) (int64, error) {
	if err := u.Conn(); err != nil {
		return 0, err
	}
	sql := "insert into " + u.table + " set nickName=?,userName=?,password=?"
	stmt, err := u.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (u *UserManager) SelectByID(userId int64) (*datamodels.User, error) {
	if err := u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "select * from " + u.table + " where id=" + strconv.FormatInt(userId, 10)
	row, err := u.mysqlConn.Query(sql)
	defer row.Close()
	if err != nil {
		return &datamodels.User{}, err
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在")
	}
	user := &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return user, nil
}
