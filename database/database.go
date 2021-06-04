package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/yuta-ron/sql-comp/config"
)

var _ DbAccessorInterface = &DbAccessor{}

type DbAccessor struct {
	Db *sql.DB
}

type DbAccessorInterface interface {
	Tables() ([]string, error)
	Columns(tblName string) (*map[string]Column, error)
}

type Column struct {
	Type   string
	IsNull bool
}

type Table struct {
	Columns *map[string]Column
}

type DBStruct struct {
	Tables map[string]Table
}

func (d *DbAccessor) Tables() ([]string, error) {
	var tbls []string
	res, err := d.Db.Query("SHOW TABLES")
	if err != nil {
		return tbls, err
	}
	defer res.Close()

	var table string
	for res.Next() {
		res.Scan(&table)
		tbls = append(tbls, table)
	}

	return tbls, nil
}

func (d *DbAccessor) Columns(tblName string) (*map[string]Column, error) {
	res, err := d.Db.Query(fmt.Sprintf("SELECT * FROM %s limit 0", tblName))
	if err != nil {
		return nil, err
	}
	defer res.Close()

	cl, _ := res.Columns()
	mCols := make(map[string]Column, len(cl))

	types, _ := res.ColumnTypes()
	for _, t := range types {
		isNull, _ := t.Nullable()
		c := &Column{
			Type:   t.DatabaseTypeName(),
			IsNull: isNull,
		}

		mCols[t.Name()] = *c
	}

	return &mCols, nil
}

func NewToDB() (*sql.DB, error) {
	toDb, err := sql.Open("mysql", config.GetToDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err = toDb.Ping(); err != nil {
		log.Fatal(err)
	}

	return toDb, nil
}

func NewFromDB() (*sql.DB, error) {
	fromDb, err := sql.Open("mysql", config.GetFromDSN())
	if err != nil {
		return nil, err
	}
	if err = fromDb.Ping(); err != nil {
		return nil, err
	}

	return fromDb, nil
}
