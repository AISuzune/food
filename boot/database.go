package boot

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	g "main/app/global"
	"os"
	"time"
)

func MysqlDBSetup() {
	config := g.Config.DataBase.Mysql

	// 创建一个新的gorm配置
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 使用彩色打印
		},
	)

	// 使用新的gorm配置打开数据库连接
	db, err := gorm.Open(mysql.Open(config.GetDsn()), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		g.Logger.Fatalf("initialize mysql db failed, err: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxIdleTime(10 * time.Second)
	sqlDB.SetConnMaxLifetime(100 * time.Second)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	err = sqlDB.Ping()
	if err != nil {
		g.Logger.Fatalf("connect to mysql db failed, err: %v", err)
	}
	g.MysqlDB = db

	g.Logger.Infof("initialize mysql db successfully")
}

func MongoDBSetup() {
	clientOptions := options.Client().ApplyURI(g.Config.DataBase.Mongo.GetAddr())
	clientOptions.SetAuth(options.Credential{
		Username: g.Config.DataBase.Mongo.Username,
		Password: g.Config.DataBase.Mongo.Password,
	})

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		g.Logger.Fatalf("initialize mongodb failed, err: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		g.Logger.Fatalf("initialize mongodb failed, err: %v", err)
	}

	g.MongoDB = client

	g.Logger.Infof("initiate mongodb successfully")
}

func RedisSetup() {
	config := g.Config.DataBase.Redis

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Addr, config.Port),
		Username: "",
		Password: config.Password,
		DB:       config.Db,
		PoolSize: 10000,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		g.Logger.Fatalf("connect to redis instance failed, err: %v", err)
	}

	g.Rdb = rdb

	g.Logger.Infof("initialize redis client successfully")
}
