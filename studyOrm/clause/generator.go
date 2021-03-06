package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
}

func genBindVar(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _values(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $tableName ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}
func _select(values ...interface{}) (string, []interface{}) {
	// VALUES ($v1), ($v2), ...
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i,value := range values{
		v := value.([]interface{})
		if bindStr == ""{
			bindStr = genBindVar(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)",bindStr))
		if i+1 != len(values){
			sql.WriteString(", ")
		}
		vars = append(vars,v...)
	}
	return sql.String(),vars
}
func _limit(values ...interface{}) (string, []interface{}) {

}
func _where(values ...interface{}) (string, []interface{}) {

}
func _orderBy(values ...interface{}) (string, []interface{}) {

}
func _insert(values ...interface{}) (string, []interface{}) {

}
