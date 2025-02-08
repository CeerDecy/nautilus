package mysql

import (
	"fmt"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Interface interface {
	DB() *gorm.DB
}

type config struct {
	Host     string `file:"host"`
	Port     string `file:"port"`
	Username string `file:"username"`
	Password string `file:"password"`
	Database string `file:"database"`
}
type provider struct {
	Cfg *config
	db  *gorm.DB
}

func (p *provider) Init(ctx servicehub.Context) error {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		p.Cfg.Username,
		p.Cfg.Password,
		p.Cfg.Host,
		p.Cfg.Port,
		p.Cfg.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	p.db = db
	logrus.Infof("connect to mysql success")
	return nil
}

func (p *provider) DB() *gorm.DB {
	return p.db
}

func init() {
	servicehub.Register("mysql-provider", &servicehub.Spec{
		Services:    []string{"mysql-provider"},
		Description: "mysql-provider service",
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
