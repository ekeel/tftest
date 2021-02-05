package model

// ValidationResult holds the results of a single field validation.
type ValidationResult struct {
	// ID is the uuid of this validation result.
	ID string
	// Type specifies the type of test scenario.
	Type string
	// Name is the name of the field that was inspected.
	Name string
	// IsMath is true if the ExpectedValue and ActualValue match.
	IsMatch bool
	// ExpectedValue is the expected value of the field.
	ExpectedValue string
	// ActualValue is the actual value of the field from AWS.
	ActualValue string
}
