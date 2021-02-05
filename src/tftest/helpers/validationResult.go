package helpers

type ValidationResult struct {
	Name string
	IsMatch bool
	ExpectedValue string
	ActualValue string
}