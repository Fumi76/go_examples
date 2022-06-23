package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

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

func print2(m *map[string]interface{}) {
	for key, val := range *m {
		fmt.Printf("%s (%v) = %v (%v)\n", key, reflect.TypeOf(key), val, reflect.TypeOf(val))
		if a, ok := val.([]interface{}); ok {
			for i, e := range a {
				fmt.Printf("%d %v (%v)\n", i, e, reflect.TypeOf(e))
				if b, ok := e.(map[string]interface{}); ok {
					print2(&b)
				}
			}
		}
	}
}

// 入力にあるキーが正当かのみ検証、値の型が一致しているかは検証しない
func validateKeys(base *map[string]interface{}, input *map[string]interface{}) []error {
	var err []error
	for key, val := range *input {

		if val2, ok := (*base)[key]; ok {

			if m, ok := val.(map[string]interface{}); ok {

				if m2, ok := val2.(map[string]interface{}); ok {

					e1 := validateKeys(&m2, &m)
					err = append(err, e1...)
				} else {
					fmt.Printf("%s 入力側がmapだが、基準側がmapではない\n", key)
				}

			} else if a, ok := val.([]interface{}); ok {
				if a2, ok := val2.([]interface{}); ok {
					// 配列の場合
					if len(a) > 0 && len(a2) > 0 {
						if m3, ok := a[0].(map[string]interface{}); ok {
							if m4, ok := a2[0].(map[string]interface{}); ok {
								e2 := validateKeys(&m4, &m3)
								err = append(err, e2...)
								continue
							}
						}
						fmt.Printf("%s どちらかが配列の要素がmapではないのでスキップ\n", key)
					} else {
						fmt.Printf("%s どちらかの配列のサイズが０なのでスキップ\n", key)
					}
				} else {
					fmt.Printf("%s 入力側が配列だが、基準側が配列ではない\n", key)
				}

			} else {
				// 入力側の値がmapではない
				fmt.Printf("%s 入力側がmapでも配列でもない %v\n", key, reflect.TypeOf(val))
			}

		} else {
			// 入力側にあるキーが基準側に無い
			err = append(err, fmt.Errorf("基準側にキー%sがない(入力が不正)", key))
			continue
		}
	}

	return err
}

/*
  yamlからa.b.cといったキーをつなげたパスを作ることは可能
*/
func main() {
	//m := make(map[interface{}]interface{})
	m := make(map[string]interface{})

	// 基準側
	var data = `
a: Easy!
b: 
  c: 2
  d: [3, 4]
e: "Text"
99: 100
f:
  - g: 5
  - h: "6"
i:
  j:
    k: {
		l: "M"
	}
o:
  - name: hoge
    value: foo
  - name: bar
    value: alpha
p:
  q:
    - name: beta
      value: delta
    - name: ceta
      value: eta    
`

	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//fmt.Printf("--- m:\n%v\n\n", m)

	print2(&m)

	// 入力側
	var j = `
	{
		"a": "Easy!",
		"b": {
			"c": 2,
			"d": [3, 4]
		}, 
		"e": "Text",
		"99": 100,
		"f": [
			{ "g": 5 },
  			{ "h": "6" }
		],
		"g": "Dummy",
		"i": {
			"j": {
				"k": {
					"n": "N"
				}
			}
		},
		"o": [
			{
				"name": "hoge",
				"value": "foo"
			},
			{
				"name": "bar",
				"value": "alpha"
			}
		],
		"p": {
			"q": [
				{
					"name1": "beta",
					"value": "delta",
					"key1": "val1"
				},
				{
					"name1": "ceta",
					"value": "eta",
					"key1": "val2"
				}
			]
		}
	}
	`

	x := map[string]interface{}{}

	err = json.Unmarshal([]byte(j), &x)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("JSON")
	print2(&x)

	// 入力にあるキーが正当かのみ検証、値の型が一致しているかは検証しない
	errs := validateKeys(&m, &x)
	for _, e := range errs {
		fmt.Printf("ERROR %v\n", e)
	}

	/*
		// キーの収集
		var keys []string
		print(m, "", &keys)
		for _, e := range keys {
			fmt.Printf("KEY %v\n", e)
		}*/
}
