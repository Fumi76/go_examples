package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/d5/tengo/v2"
)

type Hoge struct {
	tengo.ObjectImpl
	X string
	Y int
	Z *C1
}

type C1 struct {
	V1 string
}

func main() {
	// Tengo script code
	src := `
each := func(seq, fn) {
    for x in seq { fn(x) }
}

sum := 0
mul := 1
each([a, b, c, d], func(x) {
	sum += x
	mul *= x
})

s := "ほげほげ"
for a:=0; a<10; a++ {
s += a
s += e
}

s2 := ""
for k, v in m {
	s2 += k 
}

s3 := ""
s3 += hoge.X
s3 += ", " + hoge.Y

s4 := ""
s4 += hoge.Z.X
`

	// create a new Script instance
	script := tengo.NewScript([]byte(src))

	// set values
	_ = script.Add("a", 1)
	_ = script.Add("b", 9)
	_ = script.Add("c", 8)
	_ = script.Add("d", 4)
	_ = script.Add("e", "あいう")

	m := make(map[string]interface{})
	c1 := make(map[string]interface{})
	c1["key3"] = "value3"
	m["key1"] = c1
	m["key2"] = "value2"
	_ = script.Add("m", m)

	c := &C1{V1: "子要素の値１"}
	h := &Hoge{
		ObjectImpl: tengo.ObjectImpl{},
		X:          "テスト",
		Y:          100,
		Z:          c,
	}
	_ = script.Add("hoge", h)

	// run the script
	compiled, err := script.RunContext(context.Background())
	if err != nil {
		panic(err)
	}

	// retrieve values
	//sum := compiled.Get("sum")
	//mul := compiled.Get("mul")
	//s := compiled.Get("s")
	//fmt.Println(sum, mul, s) // "22 288"

	s2 := compiled.Get("s2")
	fmt.Printf("s2 %s\n", s2)

	s3 := compiled.Get("s3")
	fmt.Printf("s3 %s\n", s3)

	s4 := compiled.Get("s4")
	fmt.Printf("s4 %s\n", s4)
}

func (o *Hoge) IndexGet(index tengo.Object) (tengo.Object, error) {
	strIdx, _ := index.(*tengo.String)
	fmt.Printf("strIdx=%s\n", strIdx.Value)
	v := getField(o, strIdx.Value)
	k := reflect.ValueOf(v).Kind()
	if k == reflect.String {
		return &tengo.String{Value: v.(string)}, nil
	} else if k == reflect.Int {
		return &tengo.Int{Value: int64(v.(int))}, nil
	} else {
		return &Hoge{X: "X3"}, nil
	}
}

func getField(o *Hoge, field string) interface{} {
	r := reflect.ValueOf(o)
	f := reflect.Indirect(r).FieldByName(field)
	fmt.Printf("%+v Type=%v Kind=%v\n", f.Interface(), f.Type(), f.Kind())
	return f.Interface()
}
