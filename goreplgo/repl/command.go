package repl

// Arg represents a single argument for a REPL command, detailing its expected type and whether it is required.
type Arg struct {
	Name     string  // Name is the identifier for the argument.
	ArgType  ArgType // ArgType indicates the type of the argument (e.g., String, Int, Boolean).
	Required bool    // Required indicates whether this argument must be provided for the command to execute.
}

// Handler is a function signature for command handlers in the REPL.
// It accepts a map of argument names to their parsed values and returns a boolean and an error.
// The boolean indicates if the REPL should exit after executing this command (true means exit),
// and the error indicates any issues that occurred during execution.
type Handler func(args map[string]interface{}) (bool, error)

// Command encapsulates information about a single REPL command, including how it should be invoked
// and what it should do.
type Command struct {
	Name        string  // Name is the command's name used for invoking it.
	Description string  // Description provides a brief explanation of what the command does.
	Args        []Arg   // Args lists the arguments that the command accepts.
	Handler     Handler // Handler is the function that executes the command logic.
}
