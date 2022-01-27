package session

import (
	"testing"
	//"geeorm"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

// var (
// 	TestDB      *sql.DB
// 	TestDial, _ = dialect.GetDialect("sqlite3")
// )

// func TestMain(m *testing.M) {
// 	TestDB, _ = sql.Open("sqlite3", "../gee.db")
// 	code := m.Run()
// 	_ = TestDB.Close()
// 	os.Exit(code)
// }

// func NewSession() *Session {
// 	return New(TestDB, TestDial)
// }

func TestSession_CreateTable(t *testing.T) {
	//不需要NewEngine？？
	s := NewSession().Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}

func TestSession_Model(t *testing.T) {
	s := NewSession().Model(&User{})
	table := s.RefTable()
	s.Model(&Session{})
	if table.Name != "User" || s.RefTable().Name != "Session" {
		t.Fatal("Failed to change model")
	}
}
