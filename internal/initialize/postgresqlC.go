package initialize

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntquang/ecommerce/global"
	"go.uber.org/zap"
)

func CheckErrorPanicC(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error(err))
		panic(err)
	}
}

func InitPostgresqlC() {
	ctx := context.Background()
	m := global.Config.Postgresql
	fmt.Println("m.Host, m.Username, m.Password, m.Dbname, m.Port::", m.Host, m.Username, m.Password, m.Dbname, m.Port)
	dsn := "user=%s password=%s dbname=%s host=%s port=%d sslmode=disable"
	var s = fmt.Sprintf(dsn, m.Username, m.Password, m.Dbname, m.Host, m.Port)
	db, err := pgxpool.New(ctx, s)
	CheckErrorPanicC(err, "InitPostgresqlC initialization error")
	global.Logger.Info("initialization PostgresqlC Successfully!")

	global.Pdbc = db
}
