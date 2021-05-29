package config

import (
	"fmt"
	"os"
)

func GetFromDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("COMPARE_FROM_DB_USER"), os.Getenv("COMPARE_FROM_DB_PASSWORD"), os.Getenv("COMPARE_FROM_DB_HOST"), os.Getenv("COMPARE_FROM_DB_PORT"), os.Getenv("COMPARE_FROM_DB_NAME"))
}

func GetToDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("COMPARE_TO_DB_USER"), os.Getenv("COMPARE_TO_DB_PASSWORD"), os.Getenv("COMPARE_TO_DB_HOST"), os.Getenv("COMPARE_TO_DB_PORT"), os.Getenv("COMPARE_TO_DB_NAME"))
}
