package repl

// ArgType specifies the type of the argument expected by the command.
type ArgType int

const (
	StringArg ArgType = iota
	IntArg
	FloatArg
	BooleanArg
	DateArg
)

func (at ArgType) Name() string {
	switch at {
	case StringArg:
		return "String"
	case IntArg:
		return "Int"
	case FloatArg:
		return "Float"
	case BooleanArg:
		return "Boolean"
	case DateArg:
		return "Date"
	default:
		return "Unknown"
	}
}
