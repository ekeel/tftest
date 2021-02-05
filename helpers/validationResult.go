package helpers

type ValidationResult struct {
	ID            string
	Name          string
	IsMatch       bool
	ExpectedValue string
	ActualValue   string
}
