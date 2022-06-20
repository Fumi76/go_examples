package main

import (
	"fmt"
	"log"
	"reflect"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

var data = `
a: Easy!
b: 
  c: 2
  d: [3, 4]
e: "Text"
99: 100
f:
  - g: 5
  - h: 6
`
var kindMap = map[string]struct{}{
	"map":      {},
	"element3": {},
}

// スカラーか否かというよりは自分の下の値がキーをもっているか否かが知りたい
func isScholar(o any) bool {
	t := reflect.TypeOf(o)
	fmt.Println("Type", t, t.Kind().String())
	if _, ok := kindMap[t.Kind().String()]; ok {
		return false
	}
	return true
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func print(o any, parent string, keys *[]string) {
	if m2, ok := o.(map[interface{}]interface{}); ok {
		fmt.Printf("(%v)\n", reflect.TypeOf(m2))
		for k2, v2 := range m2 {
			fmt.Printf("%v (%v)\n", k2, reflect.TypeOf(k2))
			flag := isScholar(v2)
			if len(parent) > 0 {
				key := parent + "." + fmt.Sprintf("%v", k2)
				fmt.Printf("key %v\n", key)
				if flag {
					if slices.Contains(*keys, parent) {
						*keys = remove(*keys, parent)
					}
					*keys = append(*keys, key)
				}
				print(v2, parent+"."+fmt.Sprintf("%v", k2), keys)
			} else {
				key := fmt.Sprintf("%v", k2)
				fmt.Printf("key %v\n", key)
				if flag {
					*keys = append(*keys, key)
				}
				print(v2, fmt.Sprintf("%v", k2), keys)
			}
		}
	} else if m4, ok := o.(map[string]interface{}); ok {
		for k4, v4 := range m4 {
			fmt.Printf("%v (%v)\n", k4, reflect.TypeOf(k4))
			flag := isScholar(v4)
			if len(parent) > 0 {
				key := parent + "." + fmt.Sprintf("%v", k4)
				if flag {
					if slices.Contains(*keys, parent) {
						*keys = remove(*keys, parent)
					}
					*keys = append(*keys, key)
				}
				fmt.Printf("key %v\n", key)
				print(v4, parent+"."+fmt.Sprintf("%v", v4), keys)
			} else {
				key := fmt.Sprintf("%v", k4)
				fmt.Printf("key %v\n", key)
				if flag {
					*keys = append(*keys, key)
				}
				print(v4, fmt.Sprintf("%v", k4), keys)
			}
		}
	} else if m3, ok := o.([]interface{}); ok {
		fmt.Printf("%v (%v)\n", m3, reflect.TypeOf(m3))
		for _, e := range m3 {
			print(e, parent, keys)
		}
	} else {
		fmt.Printf("%v (%v)\n", o, reflect.TypeOf(o))
	}
}

/*
  yamlからa.b.cといったキーをつなげたパスを作ることは可能
*/
func main() {
	m := make(map[interface{}]interface{})

	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//fmt.Printf("--- m:\n%v\n\n", m)
	var keys []string
	print(m, "", &keys)
	for _, e := range keys {
		fmt.Printf("KEY %v\n", e)
	}
}
