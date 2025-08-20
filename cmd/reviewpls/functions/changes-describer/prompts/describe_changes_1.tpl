Print changes review as JSON in format, which is described as JSON-schema:
```json
{
  "type": "object",
  "properties": {
    "Summary": {
      "type": "string",
      "description": "Brief description of the changes."
    },
    "Comments": {
      "type": "array",
      "minItems": 0,
      "items": {
        "type": "object",
        "properties": {
          "Line": {
            "type": "integer",
            "description": "Start line number of the commented code."
          },
          "Text": {
            "type": "string",
            "description": "Commentary content text."
          }
        },
        "required": ["Line", "Text"],
        "additionalProperties": false
      }
    }
  },
  "required": ["Summary", "Comments"],
  "additionalProperties": false
}
```
