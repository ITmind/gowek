package repo

import (
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB
var err error

func Init() {
	initDB()
	migrate()
	SetupSuperuser()
}

func initDB() {
	slog.Info("Init db...")

	dsn := "gowek.db"
	val, ok := os.LookupEnv("DB_URL")
	if ok {
		dsn = val
	}

	DB, err = sqlx.Open("sqlite3", dsn)
	if err != nil {
		panic(err.Error())
	}

	err = DB.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func Migrate() {
	initDB()
	migrate()
}

func migrate() {
	//log.Println("Migrate...")
	// err = DB.AutoMigrate(&User{}, &Note{})
	// if err != nil {
	// 	panic(err.Error())
	// }
	//log.Println("Migrate success")
}

func SetupSuperuser() {

	superusername, ok := os.LookupEnv("SUPERUSER_USERNAME")
	if !ok {
		return
	}

	superuserpass, ok := os.LookupEnv("SUPERUSER_PASSWORD")
	if !ok {
		return
	}

	superuseremail, ok := os.LookupEnv("SUPERUSER_EMAIL")
	if !ok {
		return
	}

	//если пользоатель уже есть в базе, то не создаем
	if _, err = GetUser(superusername); err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("Setup superuser")

	superuser := User{
		Login: superusername,
		Email: superuseremail,
		Hash:  superuserpass,
	}

	err = AddUser(&superuser)
	if err != nil {
		slog.Error(err.Error())
	}
}
