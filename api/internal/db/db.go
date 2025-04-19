package db

import (
	"fmt"

	"github.com/mislu/market-api/internal/types/models"
	"github.com/mislu/market-api/internal/utils/app"
	"github.com/mislu/market-api/internal/utils/log"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(zapLogger *zap.Logger) {
	config := app.GetConfig().Database
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local&&sql_mode=ANSI_QUOTES",
		config.Username, config.Password, config.Host, config.Port, config.Database, config.Charset)
	dialector := mysql.New(mysql.Config{
		DSN: dsn,
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: log.NewGormLogger(zapLogger),
	})

	if err != nil {
		panic(err)
	}

	DB = db
	if err := autoMigrate(); err != nil {
		panic(err)
	}
}

func autoMigrate() error {
	err := DB.AutoMigrate(
		&models.User{},
	)
	return err
}
