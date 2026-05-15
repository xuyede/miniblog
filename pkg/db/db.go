package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLOptions struct {
	Host                  string        `mapstructure:"host" json:"host"`
	Username              string        `mapstructure:"username" json:"username"`
	Password              string        `mapstructure:"password" json:"password"`
	Database              string        `mapstructure:"database" json:"database"`
	MaxIdleConnections    int           `mapstructure:"max-idle-connections" json:"max-idle-connections"`
	MaxOpenConnections    int           `mapstructure:"max-open-connections" json:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `mapstructure:"max-connection-life-time" json:"max-connection-life-time"`
	LogLevel              int           `mapstructure:"log-level" json:"log-level"`
}

func (o *MySQLOptions) DSN() string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		o.Username,
		o.Password,
		o.Host,
		o.Database,
		true,
		"Local")
}

func NewMySQL(opts *MySQLOptions) (*gorm.DB, error) {
	logLevel := logger.Silent

	if opts.LogLevel != 0 {
		logLevel = logger.LogLevel(opts.LogLevel)
	}

	db, err := gorm.Open(mysql.Open(opts.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns 设置到数据库的最大打开连接数
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

	// SetConnMaxLifetime 设置连接可重用的最长时间
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	// SetMaxIdleConns 设置空闲连接池的最大连接数
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}
