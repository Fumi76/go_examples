package main

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func test1() {

	// 無くても大丈夫そう
	// "$schema": "https://json-schema.org/draft/2020-12/schema",
	// "$id": "https://example.com/product.schema.json",
	schema := `
	{

		"title": "Product",
		"description": "A product from Acme's catalog",
		"type": "object",
		"properties": {
		  "productId": {
			"description": "The unique identifier for a product",
			"type": "integer"
		  },
		  "productName": {
			"description": "Name of the product",
			"type": "string"
		  },
		  "price": {
			"description": "The price of the product",
			"type": "number",
			"exclusiveMinimum": 0
		  },
		  "tags": {
			"description": "Tags for the product",
			"type": "array",
			"items": {
			  "type": "string"
			},
			"minItems": 1,
			"uniqueItems": true
		  },
		  "dimensions": {
			"type": "object",
			"properties": {
			  "length": {
				"type": "number"
			  },
			  "width": {
				"type": "number"
			  },
			  "height": {
				"type": "number"
			  }
			},
			"required": [ "length", "width", "height" ]
		  },
		  "warehouseLocation": {
			"description": "Coordinates of the warehouse where the product is located.",
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"latitude": {
						"type": "number"
					},
					"longitude": {
						"type": "number"
					}
				},
				"required": ["latitude", "longitude"],
				"additionalProperties": false
			},
			"minItems": 2
		  }
		},
		"required": [ "productId", "productName", "price" ]
	  }
	`

	j := `
	{
		"productId": 1,
		"productName": "An ice sculpture",
		"price": 12.50,
		"tags": [ "cold", "ice" ],
		"dimensions": {
		  "length": 7.0,
		  "width": 12.0,
		  "height": 9.5
		},
		"warehouseLocation": [
			{
				"latitude": -78.75,
				"longitude": 20.4
			},
			{
				"latitude": 98.74,
				"longitude": 34.5
			}
		]
	}
	`

	schemaLoader := gojsonschema.NewStringLoader(schema)
	documentLoader := gojsonschema.NewStringLoader(j)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}

	j2 := `
	{
		"productId": 1,
		"productName": "An ice sculpture",
		"price": 12.50,
		"tags": [ "cold", "ice" ],
		"dimensions": {
		  "length": 7.0,
		  "width": 12.0,
		  "height": 9.5
		},
		"warehouseLocation": [
			{
				"latitude": -78.75,
				"longitude": 20.4
			},
			{
				"latitude": 98.74,
				"longitude": 34.5,
				"dummy": "hoge"
			}
		]
	}
	`
	// 余計なフィールドがあってもエラーにならない
	// "additionalProperties": false を指定すればチェックできる

	x := map[string]interface{}{}

	err = json.Unmarshal([]byte(j2), &x)
	if err != nil {
		fmt.Println(err)
		return
	}

	loader := gojsonschema.NewGoLoader(x)

	fmt.Println("Goのmapを使用")
	result, err = gojsonschema.Validate(schemaLoader, loader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
}

func main() {
	test1()
}
