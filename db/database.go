package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config - структура, в которой хранится информация для подключения к бд
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// Info - структура для получения данных из бд
type Info struct {
	InfoId   string `db:"info_id"`
	InfoData string `db:"info_data"`
}

// ConnectDB - функция для подключения к бд
func ConnectDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode =%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// GetDic - функция для создания кеша
func GetDic(db *sqlx.DB) (map[string][]byte, error) {
	infoArr := make([]Info, 0)
	err := db.Select(&infoArr, `select * from wb_info`)
	infoMap := make(map[string][]byte)
	for _, i := range infoArr {
		infoMap[i.InfoId] = []byte(i.InfoData)
	}
	return infoMap, err

}

// InsertDB - функция для вставки данных в бд
func InsertDB(db *sqlx.DB, id string, jsonData []byte) {
	fmt.Println(id)
	fmt.Println(string(jsonData))
	tx := db.MustBegin()
	tx.MustExec(`INSERT into wb_info (info_id,info_data) VALUES ($1,$2)`, id, jsonData)
	tx.Commit()
}
