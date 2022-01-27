package session

import (
	"database/sql"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/schema"
	"strings"
)

//负责与数据库的交互
type Session struct {
	db       *sql.DB         //数据库指针
	dialect  dialect.Dialect //go字段类型映射成数据库字段类型
	refTable *schema.Schema  //对象映射到数据库中的表
	sql      strings.Builder //sql语句
	sqlVars  []interface{}   //sql语句参数
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect}
}

//将Session的/sql语句和/sql语句参数置空
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql) //将s的内容追加到buf的缓冲区
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

//封装sql库的Exec、QueryRow、Query函数
//一是统一打印日志（包括 执行的SQL 语句和错误日志）。
//二是执行完成后，清空 (s *Session).sql 和 (s *Session).sqlVars 两个变量
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
