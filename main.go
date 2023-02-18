package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 数据库连接字符串
const dsn = "app:app@tcp(localhost:3306)/app?loc=Local&parseTime=true"

// 创建数据库链接池
func mustOpenDB(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	db := mustOpenDB(dsn)

	// 创建
	list := []TestTableName{
		{FieldOne: 1, FieldTwo: "a"},
		{FieldOne: 2, FieldTwo: "b"},
		{FieldOne: 3, FieldTwo: "c"},
	}
	err := createData(db, list)
	if err != nil {
		panic(err)
	}

	// 查询
	arg := Query{
		FieldOne: json.Number("2"),
		FieldTwo: "b",
		Current:  1,
		PageSize: 10,
	}
	data, err := fetchData(db, arg)
	if err != nil {
		panic(err)
	}

	// 打印
	for _, v := range data {
		fmt.Printf("id: %02d, fieldOne: %02d, fieldTwo: %s, createdAt: %s;\n",
			v.Id,
			v.FieldOne,
			v.FieldTwo,
			v.CreatedAt.Time.Format(time.DateTime),
		)
	}
	// id: 02, fieldOne: 02, fieldTwo: b, createdAt: 2023-02-18 17:21:06;
}

// 表实体
type TestTableName struct {
	Id        uint         // 主键
	FieldOne  int          // 字段一
	FieldTwo  string       // 字段二
	CreatedAt sql.NullTime // 创建于
}

// 创建数据
func createData(db *sql.DB, data []TestTableName) error {
	n := time.Now()

	var s strings.Builder
	s.WriteString("INSERT INTO test_table_name ")
	s.WriteString("(field_one,field_two,created_at) ")
	s.WriteString("VALUES ")

	args := make([]any, 0, len(data)*4)
	for i, v := range data {
		if i != 0 {
			s.WriteRune(',')
		}
		s.WriteString("(?,?,?)")
		args = append(args, v.FieldOne, v.FieldTwo, n)
	}

	_, err := db.Exec(s.String(), args...)
	if err != nil {
		return err
	}
	return nil
}

// 查询条件
type Query struct {
	FieldOne json.Number // 字段一
	FieldTwo string      // 字段二
	Current  int         // 当前页
	PageSize int         // 每页的数量
}

// 查询数据
func fetchData(db *sql.DB, arg Query) ([]TestTableName, error) {
	var s strings.Builder
	s.WriteString("SELECT ")
	s.WriteString("id,field_one,field_two,created_at ")
	s.WriteString("FROM test_table_name")

	var args []any
	var b bool
	if i, _ := arg.FieldOne.Int64(); i != 0 {
		b = true
		s.WriteString(" WHERE ")
		s.WriteString("field_one = ?")
		args = append(args, i)
	}
	if arg.FieldTwo != "" {
		if !b {
			b = true
			s.WriteString(" WHERE ")
		} else {
			s.WriteString(" AND ")
		}
		s.WriteString("field_two LIKE ?")
		args = append(args, "%"+arg.FieldTwo+"%")
	}

	s.WriteString(" LIMIT ? OFFSET ?")
	offset := (arg.Current - 1) * arg.PageSize
	args = append(args, arg.PageSize, offset)

	rows, err := db.Query(s.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []TestTableName
	for rows.Next() {
		var v TestTableName
		err := rows.Scan(&v.Id, &v.FieldOne, &v.FieldTwo, &v.CreatedAt)
		if err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return data, nil
}
