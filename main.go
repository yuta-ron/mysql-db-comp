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

	compare(fromInfo, toInfo)
	// compare(toInfo, fromInfo, true)

	defer fromDb.Close()
	defer toDb.Close()
}

func compare(fromInfo *database.DBStruct, toInfo *database.DBStruct) {
	fromTxt := "比較元"
	toTxt := "比較先"

	allTableNames := make(map[string]string)
	for k, _ := range fromInfo.Tables {
		allTableNames[k] = k
	}
	for k, _ := range toInfo.Tables {
		allTableNames[k] = k
	}

	for _, tblName := range allTableNames {
		if _, ok := fromInfo.Tables[tblName]; !ok {
			fmt.Printf("テーブル %s は%sのDBに存在しません\n", tblName, fromTxt)
			continue
		}
		if _, ok := toInfo.Tables[tblName]; !ok {
			fmt.Printf("テーブル %s は%sのDBに存在しません\n", tblName, toTxt)
			continue
		}

		allColumnNames := make(map[string]string)
		for clmnName, _ := range *fromInfo.Tables[tblName].Columns {
			allColumnNames[clmnName] = clmnName
		}
		for clmnName, _ := range *toInfo.Tables[tblName].Columns {
			allColumnNames[clmnName] = clmnName
		}

		for clmnName, _ := range allColumnNames {
			fromClmns := *fromInfo.Tables[tblName].Columns
			if _, ok := fromClmns[clmnName]; !ok {
				fmt.Printf("テーブル %s の %s 列は %s のテーブルに存在しません\n", tblName, clmnName, fromTxt)
				continue
			}
			toClmns := *toInfo.Tables[tblName].Columns
			if _, ok := toClmns[clmnName]; !ok {
				fmt.Printf("テーブル %s の %s 列は %s のテーブルに存在しません\n", tblName, clmnName, toTxt)
				continue
			}

			fType := fromClmns[clmnName].Type
			tType := toClmns[clmnName].Type
			if fromClmns[clmnName].Type != toClmns[clmnName].Type {
				fmt.Printf("テーブル %s の %s 列の型が異なります [比較元: %s <=> 比較先: %s]\n", tblName, clmnName, fType, tType)
			}

			fNull := fromClmns[clmnName].IsNull
			tNull := toClmns[clmnName].IsNull
			if fromClmns[clmnName].IsNull != toClmns[clmnName].IsNull {
				fmt.Printf("テーブル %s の %s 列のIsNull定義が異なります [比較元: %v <=> 比較先: %v]\n", tblName, clmnName, fNull, tNull)
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
