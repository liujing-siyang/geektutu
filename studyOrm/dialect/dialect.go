package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string                    //将 Go 语言的类型映射为数据库中的类型。
	TableExistSQL(tableName string) (string, []interface{}) //返回某个表是否存在的 SQL 语句，参数是表名(table)
}

//注册和获取 dialect 实例
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
