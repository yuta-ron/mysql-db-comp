# DBCMP

DB間のテーブル、スキーマ情報の差分を検出するツール

### installation

```
git clone git@github.com:yuta-ron/mysql-db-comp.git
make build-linux-amd64
mv dbcmp <パスが通っているところ>
```
### Setup
下記の環境変数が設定されていないと動作しません。

Example: 
```
export COMPARE_FROM_DB_HOST=127.0.0.1
export COMPARE_FROM_DB_PORT=3306
export COMPARE_FROM_DB_USER=root
export COMPARE_FROM_DB_PASSWORD=root
export COMPARE_FROM_DB_NAME=mydb

export COMPARE_TO_DB_HOST=127.0.0.1
export COMPARE_TO_DB_PORT=3307
export COMPARE_TO_DB_USER=root
export COMPARE_TO_DB_PASSWORD=root
export COMPARE_TO_DB_NAME=mydb
```

### Options

```
dbcmp help
```
