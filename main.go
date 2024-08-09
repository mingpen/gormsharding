package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/sharding"
)

const (
	dsn = "postgres://storagetest:storagetest@127.0.0.1:5432/storagetest?sslmode=disable"
)

var (
	db *gorm.DB
)

func main() {
	initDB()

	// 使用分表
	err := db.Use(sharding.Register(sharding.Config{
		ShardingKey:         "file_id", // 分表的字段
		NumberOfShards:      128,       // 多少张子表
		PrimaryKeyGenerator: sharding.PKPGSequence,
		// PrimaryKeyGenerator: sharding.PKSnowflake,
	}, &FileHisRecord{}))
	if err != nil {
		panic(err)
	}

	// 自动创建表（也可以 db.Exec(`CREATE TABLE ` ...... ）
	err = db.AutoMigrate(&FileHisRecord{})
	if err != nil {
		panic(err)
	}

	// 清理数据
	// result := db.Delete(&FileHisRecord{}, &FileHisRecord{FileId: "file1"})
	result := db.Delete(&FileHisRecord{}, "file_id = ?", "file1")
	if result.Error != nil {
		panic(result.Error)
	}
	if result.RowsAffected > 0 {
		println("cleaned")
	}

	// 插入数据
	err = db.Create(&FileHisRecord{FileId: "file1", Version: 1}).Error
	if err != nil {
		panic(err)
	}

	// 查询数据
	var record FileHisRecord
	// err = db.First(&record, &FileHisRecord{FileId: "file1"}).Error
	err = db.First(&record, "file_id = ?", "file1").Error
	if err != nil {
		panic(err)
	}
	log.Printf("created record:"+logger.Red+" %+v\n"+logger.Reset, record)

	// 更新数据
	err = db.Model(&FileHisRecord{}).Where("file_id = ?", "file1").Update("version", 2).Error
	if err != nil {
		panic(err)
	}
	// 查询数据
	var newRecord FileHisRecord
	err = db.First(&newRecord, "file_id = ?", "file1").Error
	if err != nil {
		panic(err)
	}
	log.Printf("updated record:"+logger.Red+" %+v\n"+logger.Reset, newRecord)
}

func initDB() {
	var err error
	db, err = gorm.Open(postgres.New(postgres.Config{DSN: dsn}),
		&gorm.Config{Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		})})
	if err != nil {
		panic(err)
	}
}
