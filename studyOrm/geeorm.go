package geeorm

import (
	"database/sql"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/session"
)

//与用户交互
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, soucer string) (e *Engine, err error) {
	//打开由数据库驱动程序名称和驱动程序特定数据源名称指定的数据库
	db, err := sql.Open(driver, soucer)
	if err != nil {
		log.Error(err)
		return
	}
	//验证到数据库的连接是否仍然存在，并在必要时建立连接
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{
		db:      db,
		dialect: dial,
	}
	log.Info("Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}
