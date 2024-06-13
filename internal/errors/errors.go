package errors

import "fmt"

type NotFoundError struct {
	StockID string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Stock ID %s not found.", e.StockID)
}

type InvalidInputError struct {
	Param string
	Value string
}

func (e *InvalidInputError) Error() string {
	return fmt.Sprintf("Invalid value for %s: %s.", e.Param, e.Value)
}

type ExternalAPIError struct {
	Message string
}

func (e *ExternalAPIError) Error() string {
	return fmt.Sprintf("An external API error has occured. %s", e.Message)
}
