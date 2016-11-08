package check

import "fmt"

// BaseCheck is modelling a basic check.
type BaseCheck struct {
	id         string
	name       string
	comment    string
	expected   interface{}
	lastResult Result
}

// GetID returns ID
func (b *BaseCheck) GetID() string {
	return b.id
}

// SetID sets ID
func (b *BaseCheck) SetID(id string) {
	b.id = id
}

// GetName returns name
func (b *BaseCheck) GetName() string {
	return b.name
}

// SetName sets name
func (b *BaseCheck) SetName(name string) {
	b.name = name
}

// GetComment returns comment
func (b *BaseCheck) GetComment() string {
	return b.comment
}

// SetComment sets comment
func (b *BaseCheck) SetComment(comment string) {
	b.comment = comment
}

// GetExpected returns the expected value of the check
func (b *BaseCheck) GetExpected() interface{} {
	return b.expected
}

// SetComment sets the exptected value of the check
func (b *BaseCheck) SetExpected(expected interface{}) {
	b.expected = expected
}

// GetStatus returns current status
func (b *BaseCheck) GetStatus() Status {
	return b.lastResult.Status
}

// SetResult stores a result.
func (b *BaseCheck) SetResult(result Result) {
	b.lastResult = result
}

// GetStatusString returns current status
func (b *BaseCheck) GetStatusString() string {
	return b.GetLastResult().GetStatusString()
}

// GetLastResult returns full last result
func (b *BaseCheck) GetLastResult() Result {
	return b.lastResult
}

// Execute executes this check
func (b *BaseCheck) Execute() Result {
	// implement in subclasses
	return Result{
		Status: StatusOK,
	}
}

// ParseBaseData parses basic check config
func (b *BaseCheck) ParseBaseData(configMap ConfigMap) {
	if b.GetID() == "" {
		id, err := configMap.GetString("id")
		if err == nil {
			b.SetID(id)
		}
	}

	if b.GetName() == "" {
		name, err := configMap.GetString("name")
		if err == nil {
			b.SetName(name)
		}
	}

	if b.GetComment() == "" {
		comment, err := configMap.GetString("comment")
		if err == nil {
			b.SetComment(comment)
		}
	}

	if b.GetExpected() == nil {
		expected, err := configMap.GetString("expected")
		if err == nil {
			b.SetExpected(expected)
		}
	}
}

func (b *BaseCheck) Parse(configMap ConfigMap) {
	b.ParseBaseData(configMap)
}

// String formats
func (b *BaseCheck) String() string {
	message := ""
	if b.GetStatus() != StatusOK && b.GetStatus() != StatusUnknown {
		message = " (" + b.GetLastResult().Message + ")"
	}
	id := ""
	if b.GetID() != "" {
		message = " [" + b.GetID() + "]"
	}
	return fmt.Sprintf("%s - %s%s%s", b.GetStatusString(), b.GetName(), id, message)
}

func (b *BaseCheck) getEmptyResult() Result {
	return Result{
		Status:   StatusUnknown,
		Expected: b.GetExpected(),
	}
}
