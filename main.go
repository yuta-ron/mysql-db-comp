package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yuta-ron/sql-comp/database"
)

func main() {
	// DB initialization
	fromDb, err := database.NewFromDB()
	if err != nil {
		log.Fatal(err)
	}
	toDb, err := database.NewToDB()
	if err != nil {
		log.Fatal(err)
	}

	fAccr := &database.DbAccessor{Db: fromDb}
	tAccr := &database.DbAccessor{Db: toDb}

	fromInfo, err := makeInfo(fAccr)
	if err != nil {
		log.Fatal(err)
	}
	toInfo, err := makeInfo(tAccr)
	if err != nil {
		log.Fatal(err)
	}

	compare(fromInfo, toInfo, false)
	compare(toInfo, fromInfo, true)

	defer fromDb.Close()
	defer toDb.Close()
}

func compare(fromInfo *database.DBStruct, toInfo *database.DBStruct, reversed bool) {
	fromTxt := "比較元"
	toTxt := "比較先"
	if reversed {
		fromTxt = "比較先"
		toTxt = "比較元"
	}

	for tblName, t := range fromInfo.Tables {
		if _, ok := toInfo.Tables[tblName]; !ok {
			fmt.Printf("テーブル %s は%sのDBに存在しません", tblName, fromTxt)
			continue
		}

		toClmn := *toInfo.Tables[tblName].Columns

		for clmnName, fromColumn := range *t.Columns {
			if _, ok := toClmn[clmnName]; !ok {
				fmt.Printf("列名 %s は%sのDBに存在しません", clmnName, fromTxt)
				continue
			}

			toColumn := toClmn[clmnName]
			if fromColumn.Type != toColumn.Type {
				fmt.Printf("テーブル名: %s 列名 %s の型定義が異なります。%s=%v , %s=%v", tblName, clmnName, fromTxt, fromColumn.Type, toTxt, toColumn.Type)
			}
			if fromColumn.Length != toColumn.Length {
				fmt.Printf("テーブル名: %s 列名 %s のLengthが異なります。%s=%v , %s=%v", tblName, clmnName, fromTxt, fromColumn.Length, toTxt, toColumn.Length)
			}
			if fromColumn.IsNull != toColumn.IsNull {
				fmt.Printf("テーブル名: %s 列名 %s のIsNull定義が異なります。%s=%v , %s=%v", tblName, clmnName, fromTxt, fromColumn.IsNull, toTxt, toColumn.IsNull)
			}
		}
	}
}

func makeInfo(acc *database.DbAccessor) (*database.DBStruct, error) {
	ds := &database.DBStruct{
		Tables: make(map[string]database.Table),
	}

	tbls, err := acc.Tables()
	if err != nil {
		return nil, err
	}

	for _, tableName := range tbls {
		columns, err := acc.Columns(tableName)
		if err != nil {
			return nil, err
		}

		ds.Tables[tableName] = database.Table{
			Columns: columns,
		}
	}

	return ds, nil
}
