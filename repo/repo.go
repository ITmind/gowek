package repo

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"

	sl3 "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func Init() {
	initDB()
	runMigration()
	SetupSuperuser()
}

func getDBName() string {

	dsn := "gowek.db"
	val, ok := os.LookupEnv("DB_URL")
	if ok {
		dsn = val
	}

	return dsn
}

func initDB() {
	slog.Info("Init db...")

	dsn := getDBName()

	var err error
	db, err = sqlx.Open("sqlite3", dsn)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func Migrate() {
	initDB()
	runMigration()
}

func Close() {
	slog.Info("Close db...")
	err := db.DB.Close()
	if err != nil {
		slog.Error(err.Error())
	}
}

func runMigration() {
	slog.Info("Migrate...")

	instance, err := sl3.WithInstance(db.DB, &sl3.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "sqlite3", instance)
	//dsn := getDBName()
	//m, err := migrate.New("file://migrations", "sqlite3://"+dsn)
	if err != nil {
		slog.Error(err.Error())
		panic(err.Error())
	}
	m.Up()

	if err != nil {
		slog.Error(err.Error())
		panic(err.Error())
	}
	slog.Info("Migrate success")
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
	user, err := GetUser(superusername)
	if err != nil && err != sql.ErrNoRows {
		panic("Error create superuser: " + err.Error())
	}

	var eq = user == (User{})
	if !eq {
		return
	}

	slog.Info("Setup superuser")

	superuser := User{
		Login:   superusername,
		Email:   superuseremail,
		Hash:    superuserpass,
		IsAdmin: true,
	}

	err = AddUser(&superuser)
	if err != nil {
		panic("Error create superuser: " + err.Error())
	}
}
