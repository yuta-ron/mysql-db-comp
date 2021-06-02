package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yuta-ron/sql-comp/database"
)

const (
	codeOk    = 0
	codeError = 1
)

var green func(a ...interface{}) string
var yellow func(a ...interface{}) string
var red func(a ...interface{}) string

func init() {
	green = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	red = color.New(color.FgRed).SprintFunc()
}

func main() {
	os.Exit(execute())
}

func execute() int {
	flag.Parse()

	if len(flag.Args()) == 0 {
		return run()
	}
	arg := flag.Args()[0]
	if strings.Contains(arg, "run") {
		return run()
	}
	if strings.Contains(arg, "doctor") {
		return doctor()
	}
	if strings.Contains(arg, "help") {
		return showHelp()
	}
	if strings.Contains(arg, "info") {
		return info()
	}

	fmt.Printf("%s option is not found\n", flag.Args())
	return codeError
}

func run() int {
	if !environmentCheck() {
		fmt.Println("環境変数が設定されていません")

		return codeError
	}

	// DB initialization
	fromDb, err := database.NewFromDB()
	if err != nil {
		log.Fatal(err)
	}
	defer fromDb.Close()

	toDb, err := database.NewToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer toDb.Close()

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

	if compare(fromInfo, toInfo) {
		return codeOk
	}

	return codeError
}

func compare(fromInfo *database.DBStruct, toInfo *database.DBStruct) bool {
	result := true

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
			log := fmt.Sprintf("テーブル %s は%sのDBに存在しません\n", tblName, fromTxt)
			fmt.Println(red(log))
			continue
		}
		if _, ok := toInfo.Tables[tblName]; !ok {
			log := fmt.Sprintf("テーブル %s は%sのDBに存在しません\n", tblName, toTxt)
			fmt.Println(red(log))
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
				log := fmt.Sprintf("テーブル %s の %s 列は %s のテーブルに存在しません\n", tblName, clmnName, fromTxt)
				fmt.Println(red(log))
				result = false
				continue
			}
			toClmns := *toInfo.Tables[tblName].Columns
			if _, ok := toClmns[clmnName]; !ok {
				log := fmt.Sprintf("テーブル %s の %s 列は %s のテーブルに存在しません\n", tblName, clmnName, toTxt)
				fmt.Println(red(log))
				result = false
				continue
			}

			fType := fromClmns[clmnName].Type
			tType := toClmns[clmnName].Type
			if fromClmns[clmnName].Type != toClmns[clmnName].Type {
				log := fmt.Sprintf("テーブル %s の %s 列の型が異なります [比較元: %s <=> 比較先: %s]\n", tblName, clmnName, fType, tType)
				fmt.Println(red(log))
				result = false
			}

			fNull := fromClmns[clmnName].IsNull
			tNull := toClmns[clmnName].IsNull
			if fromClmns[clmnName].IsNull != toClmns[clmnName].IsNull {
				log := fmt.Sprintf("テーブル %s の %s 列のIsNull定義が異なります [比較元: %v <=> 比較先: %v]\n", tblName, clmnName, fNull, tNull)
				fmt.Println(yellow(log))
				result = false
			}
		}
	}

	if result {
		fmt.Println(green("DBは同期されています"))
	}

	return result
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

func info() int {
	logo :=
		`
____  ____   ____ __  __ ____
|  _ \| __ ) / ___|  \/  |  _ \
| | | |  _ \| |   | |\/| | |_) |
| |_| | |_) | |___| |  | |  __/
|____/|____/ \____|_|  |_|_|
`
	fmt.Println(logo)

	fmt.Println("Please contact @yutaro-nishi if needs more info.")

	return codeOk
}

func showHelp() int {
	help :=
		`
doctor: 環境変数のチェックを実行します。
run: 比較を実行します
help: コマンド一覧を表示します
info: 情報を表示します
`
	fmt.Println(help)

	return codeOk
}

func doctor() int {
	if environmentCheck() {
		fmt.Println(green("環境変数は正常に設定されています"))
		return codeOk
	}

	fmt.Println("環境変数が設定されていません")
	return codeError
}

func environmentCheck() bool {
	ok := true

	_, exists := os.LookupEnv("COMPARE_FROM_DB_HOST")
	ok = ok && exists
	_, exists = os.LookupEnv("COMPARE_FROM_DB_PORT")
	ok = ok && exists
	_, exists = os.LookupEnv("COMPARE_FROM_DB_USER")
	ok = ok && exists
	_, exists = os.LookupEnv("COMPARE_FROM_DB_PASSWORD")
	ok = ok && exists
	_, exists = os.LookupEnv("COMPARE_FROM_DB_NAME")
	ok = ok && exists
	_, exists = os.LookupEnv("COMPARE_TO_DB_HOST")
	ok = ok && exists
	_, exists = os.LookupEnv("COMPARE_TO_DB_PORT")
	ok = ok && exists
	_, exists = os.LookupEnv("COMPARE_TO_DB_USER")
	ok = ok && exists
	_, exists = os.LookupEnv("COMPARE_TO_DB_PASSWORD")
	ok = ok && exists
	_, exists = os.LookupEnv("COMPARE_TO_DB_NAME")
	ok = ok && exists

	return ok
}
