package check

import "fmt"

// Status is the result status of a check
type Status int

const (
	StatusUnknown Status = iota
	StatusOK
	StatusNotice
	StatusWarning
	StatusError
	StatusCritical
)

// Result is the result of a check (status plus optional message and detail info)
type Result struct {
	Status   Status
	Value    interface{}
	Expected interface{}
	Message  string
	Info     string
}

// GetStatusString returns a printable/reable string for the result's status
func (r Result) GetStatusString() string {
	switch r.Status {
	case StatusUnknown:
		return "Unknown"
	case StatusOK:
		return "OK"
	case StatusNotice:
		return "Notice"
	case StatusWarning:
		return "Warning"
	case StatusError:
		return "Error"
	case StatusCritical:
		return "Critical"
	}
	return ""
}

func (r Result) String() string {
	return fmt.Sprintf("%s (Value: %s, Expected: %s, Message: %s, Info: %s)", r.GetStatusString(), r.Value, r.Expected, r.Message, r.Info)
}
