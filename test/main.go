package main

import (
	"fmt"
	"reflect"
	"encoding/json"
	"log"
)

type Student struct {
	Name  string
	Age   int
	Score float32
}

type Student2 struct {
	Name  string
	Age   int
	Score float32
}

func test(b interface{}) {
	t := reflect.TypeOf(b)
	fmt.Println(t)

	v := reflect.ValueOf(b)
	fmt.Println(v)

	k := v.Kind()
	fmt.Println(k)


	iv := v.Interface()
	fmt.Println(iv)


	stu, ok := iv.(Student)
	if ok {
		fmt.Printf("%v %T\n", stu, stu)
	}
	stu2 := &Student2{}
	jsons, _ := json.Marshal(b) //转换成JSON返回的是byte[]
	json.Unmarshal(jsons,stu2)
	log.Println(stu2)
}

func main() {
	var a Student = Student{
		Name:  "stu01",
		Age:   18,
		Score: 92,
	}
	//test(a)
	stu2 := &Student2{}
	jsons, _ := json.Marshal(a) //转换成JSON返回的是byte[]
	json.Unmarshal(jsons,stu2)
	log.Println(stu2)
}
