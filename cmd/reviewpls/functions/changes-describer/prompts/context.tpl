<context>
The following items were attached by the user. They are up-to-date and don't need to be re-read.

<user_rules>
The user has specified the following rules that should be applied:

Rules title: File line numbers
````
The user can provide data with line numbers added in a text file, for example:
```go
1:package main
2:
3:import (
4:	"fmt"
5:)
6:
7:func main() {
8:	fmt.Println("Bye World")
9:	fmt.Println("Exit")
10:}
```

Example dialogue using this data:
User:
Print line 8.
Assistant:
```go
	fmt.Println("Bye World")
```
````

Rules title: Description of patch format
````
Stream patch format that contains only the final result:
```blame_patch
<|file|>relative/path/to/the/file.txt
<|commits|>
<|commit|>019bc34<|commit_message|>Commit message text for the commit 019bc34
<|commit|>e4cebe2<|commit_message|>Commit message text for the commit e4cebe2
<|patch_content|>
<|equal|>
1:Equal content.
<|commit|>019bc34<|add|>
2:Added content in commit 019bc34
<|commit|>e4cebe2<|add|>
3:Added content in commit e4cebe2
```

Example of the original file `main.go`:
```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World")
}
```

Example of the file with the patch applied to `main.go`:
```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Bye World")
	fmt.Println("Exit")
}
```

Patched file `main.go` with added line numbers (line numbers are not part of the program's source code, they are added for navigation convenience):
```go
1:package main
2:
3:import (
4:	"fmt"
5:)
6:
7:func main() {
8:	fmt.Println("Bye World")
9:	fmt.Println("Exit")
10:}
```

Example of the patch body applied to `main.go`:
```blame_patch
<|file|>main.go
<|commits|>
<|commit|>b053442<|commit_message|>Change printed message
<|patch_content|>
<|equal|>
1:package main
2:
3:import (
4:	"fmt"
5:)
6:
7:func main() {
<|commit|>b053442<|add>
8:	fmt.Println("Bye World")
9:	fmt.Println("Exit")
10:}
````

Rules title: Following JSON-schema
````
## Example dialog
User:
<example_json_schema>
```json
{
	"type": "object",
	"properties": {
		"steps": {
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"explanation": {
						"type": "string"
					},
					"output": {
						"type": "string"
					}
				},
				"required": [
					"explanation",
					"output"
				],
				"additionalProperties": false
			}
		},
		"final_answer": {
			"type": "string"
		}
	},
	"required": [
		"steps",
		"final_answer"
	],
	"additionalProperties": false
}
</example_json_schema>
```

Follow <example_json_schema> to answer in JSON on this question: how can I solve 8x + 7 = -23?

Assistant:
{
	"steps": [
		{
			"explanation": "Start with the equation 8x + 7 = -23.",
			"output": "8x + 7 = -23"
		},
		{
			"explanation": "Subtract 7 from both sides to isolate the term with the variable.",
			"output": "8x = -23 - 7"
		},
		{
			"explanation": "Simplify the right side of the equation.",
			"output": "8x = -30"
		},
		{
			"explanation": "Divide both sides by 8 to solve for x.",
			"output": "x = -30 / 8"
		},
		{
			"explanation": "Simplify the fraction.",
			"output": "x = -15 / 4"
		}
	],
	"final_answer": "x = -15 / 4"
}

````
</user_rules>
</context>