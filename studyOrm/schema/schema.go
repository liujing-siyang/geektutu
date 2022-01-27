package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

type Field struct {
	Name string //字段名
	Type string //字段类型
	Tag  string //字段标签
}

//对象(object)和表(table)的转换
type Schema struct {
	Model      interface{} //映射的对象
	Name       string      //表名
	Fields     []*Field    //字段
	FieldNames []string
	fieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()//设计的入参dest是一个对象(表结构体)的指针，因此需要 reflect.Indirect() 获取指针指向的实例
	schema := Schema{
		Model: dest,
		Name: modelType.Name(),//结构体的名称作为表名
		fieldMap: map[string]*Field{},
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name){//非匿名字段并以大写字母开头
			field := Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
				//Tag: string(p.Tag),
			}
			if v,ok := p.Tag.Lookup("geeorm");ok{
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, &field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = &field
		}
	}

	return &schema
}
