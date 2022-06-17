package main

import (
	"fmt"
	"reflect"
)

type Hoge struct {
	X string
	Y int
	Z *C1
}

type C1 struct {
	V1 string
}

func main() {

	h := &Hoge{
		X: "bbb",
		Y: 100,
		Z: &C1{
			V1: "aaa",
		},
	}

	v := reflect.ValueOf(h)
	fmt.Println(v.Type(), v.Kind()) // *main.Hoge ptr
	e := v.Elem()
	fmt.Println(e.Type(), e.Kind()) // main.Hoge struct
	x := v.Elem().FieldByName("X")
	fmt.Println(x.String()) // bbb
	z := v.Elem().FieldByName("Z")
	fmt.Println(z.Type(), z.Kind()) // *main.C1 ptr
	ze := z.Elem()
	fmt.Println(ze.Type(), ze.Kind()) // main.C1 struct
	v1 := z.Elem().FieldByName("V1")
	fmt.Println(v1.String()) // aaa

	/*
		t := reflect.TypeOf(h)
		//v := reflect.ValueOf(h)
		fmt.Println(t)
		fmt.Println(t.Kind()) // ptr
		fmt.Println(t.Elem()) // main.Hoge
		sf, found := t.Elem().FieldByName("Z")
		if found {
			fmt.Println("found")
			ft := reflect.TypeOf(sf)
			fmt.Println(ft.Kind())
			f, found := ft.Elem().FieldByName("V1")
			if found {
				fmt.Println(f.Type.Kind())
			} else {
				fmt.Println("Not found V1")
			}
		}
	*/
}
