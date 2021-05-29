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

	for tblName, t := range fromInfo.Tables {
		if _, ok := toInfo.Tables[tblName]; !ok {
			fmt.Printf("テーブル %s は比較先のDBに存在しません", tblName)
			continue
		}

		toClmn := *toInfo.Tables[tblName].Columns

		for clmnName, fromColumn := range *t.Columns {
			fmt.Println(clmnName)
			if _, ok := toClmn[clmnName]; !ok {
				fmt.Printf("フィールド %s は比較先のDBに存在しません", clmnName)
				continue
			}

			toColumn := toClmn[clmnName]
			if fromColumn.Type != toColumn.Type {
				fmt.Printf("フィールド %s の型定義が異なります。比較元=%v , 比較先=%v", clmnName, fromColumn.Type, toColumn.Type)
			}
			if fromColumn.Length != toColumn.Length {
				fmt.Printf("フィールド %s のLengthが異なります。比較元=%v , 比較先=%v", clmnName, fromColumn.Length, toColumn.Length)
			}
			if fromColumn.IsNull != toColumn.IsNull {
				fmt.Printf("フィールド %s のIsNull定義が異なります。比較元=%v , 比較先=%v", clmnName, fromColumn.IsNull, toColumn.IsNull)
			}
		}
	}

	defer fromDb.Close()
	defer toDb.Close()
}

func compare(fromInfo *database.DBStruct, toInfo *database.DBStruct) error {

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
