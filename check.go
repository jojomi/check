package check

// Check is modelling a check to be executed
type Check interface {
	GetID() string
	SetID(id string)

	GetName() string
	SetName(name string)

	GetComment() string
	SetComment(comment string)

	GetStatus() Status
	GetStatusString() string
	GetLastResult() Result
	SetResult(result Result)

	GetExpected() interface{}
	SetExpected(value interface{})

	Execute() Result
	Parse(configMap ConfigMap)
	String() string
}
