package schema_test

import (
	"encoding/json"
	"fmt"

	. "github.com/WinPooh32/reviewpls/internal/schema"
)

func Example() {
	object := Object().
		Title("object title").
		Description("object description").
		Nullable(false).
		AdditionalProperties(false).
		Required(
			"field_str",
			"field_int",
		).
		Property("field_str", String().
			Nullable(true).
			Description("My super-duper string field description").
			Enum("one", "two", "three").EnumNull(),
		).
		Property("field_int",
			Integer().
				Enum(1, 2, 3),
		).
		Property("field_number",
			Number().
				Enum(0.1, 2, 3),
		).
		Property("field_number_array",
			Array().
				Description("list of numbers").
				Items(Integer()).
				MinItems(10).
				MaxItems(20),
		)

	bs, err := json.MarshalIndent(&object, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bs))
	// Output: {
	// 	"type": "object",
	// 	"title": "object title",
	// 	"description": "object description",
	// 	"nullable": false,
	// 	"additionalProperties": false,
	// 	"required": [
	// 		"field_str",
	// 		"field_int"
	// 	],
	// 	"properties": {
	// 		"field_int": {
	// 			"type": "integer",
	// 			"enum": [
	// 				1,
	// 				2,
	// 				3
	// 			]
	// 		},
	// 		"field_number": {
	// 			"type": "number",
	// 			"enum": [
	// 				0.1,
	// 				2,
	// 				3
	// 			]
	// 		},
	// 		"field_number_array": {
	// 			"type": "array",
	// 			"description": "list of numbers",
	// 			"items": {
	// 				"type": "integer"
	// 			},
	// 			"minItems": 10,
	// 			"maxItems": 20
	// 		},
	// 		"field_str": {
	// 			"type": "string",
	// 			"description": "My super-duper string field description",
	// 			"nullable": true,
	// 			"enum": [
	// 				"one",
	// 				"two",
	// 				"three",
	// 				null
	// 			]
	// 		}
	// 	}
	// }
}
