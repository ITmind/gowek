package repo

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB
var err error

func Init() {
	//Обычная функция, т.к. нам нужно сначала переменные среды загрузить
	initDB()
	migrate()
	SetupSuperuser()
}

func initDB() {
	log.Println("Init db...")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	dsn := "gorm.db"
	val, ok := os.LookupEnv("DB_URL")
	if ok {
		dsn = val
	}

	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
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
	err = DB.AutoMigrate(&User{}, &Note{})
	if err != nil {
		panic(err.Error())
	}
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

	var count int64
	if err := DB.Model(&User{}).Where("login=?", superusername).Count(&count).Error; err != nil {
		log.Println(err)
		return
	}

	if count != 0 {
		return
	}

	log.Println("Setup superuser")

	superuser := &User{
		Login: superusername,
		Email: superuseremail,
		Hash:  superuserpass,
	}

	DB.Create(superuser)
}

func CheckUser(login, hash string) bool {
	var count int64
	DB.Model(&User{}).Where("login=? and Hash=?", login, hash).Count(&count)

	return count != 0
}

func GetUser(login, hash string) (User, bool) {
	var user User
	res := DB.Select("ID", "login").First(&user, "login=? and Hash=?", login, hash)
	if res.Error != nil {
		log.Println(res.Error)
		return User{}, false
	}

	return user, true
}

func GetUserID(login string) (uint, error) {
	var user User
	res := DB.Select("ID", "login").First(&user, "login=?", login)
	return user.ID, res.Error
}

func GetEntities[T User | Note]() []T {
	var list []T
	DB.Find(&list)
	return list
}

func GetEntity[T User | Note](id uint) T {
	var obj T
	DB.First(&obj, "id = ?", id)
	return obj
}

func GetEntitiesByUser[T Note](userID uint) []T {
	var entities []T
	DB.Find(&entities, "user_id = ?", userID)
	return entities
}

func AddEntity[T User | Note](obj *T) {
	DB.Create(&obj)
}
