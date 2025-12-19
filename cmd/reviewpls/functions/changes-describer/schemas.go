package changesdescriber

import (
	sc "github.com/WinPooh32/reviewpls/internal/schema"
)

var responseChangesSchema = sc.Object().
	AdditionalProperties(false).
	Required(
		"Summary",
		"Hypotheses",
	).
	Property("Summary",
		sc.String().
			Description("Brief description of the changes."),
	).
	Property("Hypotheses",
		sc.Array().
			MinItems(0).
			Items(
				sc.Object().
					AdditionalProperties(false).
					Required(
						"Line",
						"Hypothesis",
					).
					Property("Line",
						sc.Integer().
							Description("Start line number of the commented code."),
					).
					Property("Hypothesis",
						sc.String().
							Description("Hypothesis content text."),
					),
			),
	).
	JSON()
