package geecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var dbtest = map[string][]byte{
	"Tom":  []byte("630"),
	"Jack": []byte("589"),
	"Sam":  []byte("567"),
}

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatal("callback failed")
	}
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(dbtest))
	gee := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := dbtest[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key]++
				return v, nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range dbtest {
		if view, err := gee.Get(k); err == nil || view.String() == string(v) {
			str := "1000"
			view.b = []byte(str) //getLocally中，value := ByteView{b: bytes}，切片没有copy,赋值操作结果没有影响到db???
			fmt.Println("sucess")
		}
		if view, err := gee.Get(k); err != nil || view.String() != string(v) {
			t.Fatal("failed to get value of Tom")
		}
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	for k, v := range dbtest {
		fmt.Println(k, string(v))
	}

	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}

// func TestGetGroup(t *testing.T) {
// 	groupName := "scores"
// 	NewGroup(groupName, 2<<10, GetterFunc(
// 		func(key string) (bytes []byte, err error) { return }))
// 	if group := GetGroup(groupName); group == nil || group.name != groupName {
// 		t.Fatalf("group %s not exist", groupName)
// 	}

// 	if group := GetGroup(groupName + "111"); group != nil {
// 		t.Fatalf("expect nil, but %s got", group.name)
// 	}
// }
