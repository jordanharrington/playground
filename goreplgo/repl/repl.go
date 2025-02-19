package repl

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Prompt defines a function type that returns a string to be used as the REPL prompt.
type Prompt func() string

// REPL represents the structure for a read-eval-print loop (REPL) containing commands.
type REPL struct {
	commands map[string]*Command // Map of command names to Command objects
	prompt   Prompt              // Function to generate the prompt string
}

// NewREPL creates a new REPL instance with a default prompt.
func NewREPL() *REPL {
	return WithPrompt(func() string {
		return ">> " // Default prompt
	})
}

// WithPrompt creates a new REPL instance with a specified prompt function.
func WithPrompt(prompt Prompt) *REPL {
	return &REPL{
		commands: make(map[string]*Command),
		prompt:   prompt,
	}
}

// Register adds a new command to the REPL with specified properties.
func (r *REPL) Register(name, description string, args []Arg, handler Handler) *REPL {
	r.commands[name] = &Command{
		Name:        name,
		Description: description,
		Args:        args,
		Handler:     handler,
	}
	return r
}

// Add adds a pre-defined Command object to the REPL.
func (r *REPL) Add(command *Command) *REPL {
	r.commands[command.Name] = command
	return r
}

// Run starts the REPL, continuously accepting input and executing commands.
func (r *REPL) Run() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(r.prompt())
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		tokens := parseInputWithQuotes(input)
		commandName := tokens[0]
		command, exists := r.commands[commandName]
		if !exists {
			fmt.Println("Unknown command:", commandName)
			continue
		}

		if len(tokens)-1 < len(command.Args) {
			fmt.Println("Not enough arguments. Usage:")
			r.printCommandUsage(commandName)
			continue
		}

		argValues := make(map[string]interface{})
		for i, arg := range command.Args {
			rawValue := tokens[i+1]
			var err error
			var value interface{}

			switch arg.ArgType {
			case IntArg:
				value, err = strconv.Atoi(rawValue)
				if err != nil {
					fmt.Printf("Invalid value for '%s': %s\n", arg.Name, err)
					continue
				}
			case StringArg:
				value = rawValue
			case FloatArg:
				value, err = strconv.ParseFloat(rawValue, 64)
				if err != nil {
					fmt.Printf("Invalid value for '%s': %s\n", arg.Name, err)
					continue
				}
			case BooleanArg:
				value, err = strconv.ParseBool(rawValue)
				if err != nil {
					fmt.Printf("Invalid value for '%s': %s\n", arg.Name, err)
					continue
				}
			case DateArg:
				value, err = time.Parse("2006-01-02", rawValue)
				if err != nil {
					fmt.Printf("Invalid date format for '%s' (expected YYYY-MM-DD): %s\n", arg.Name, err)
					continue
				}
			default:
				fmt.Println("Unsupported argument type.")
				continue
			}

			argValues[arg.Name] = value
		}

		exit, err := command.Handler(argValues)
		if err != nil {
			fmt.Printf("Error executing command: %s\n", err)
		}

		if exit {
			fmt.Printf("Exitting...")
			return
		}
	}
}

func parseInputWithQuotes(input string) []string {
	var tokens []string
	inQuotes := false
	currentToken := strings.Builder{}

	for _, r := range input {
		switch {
		case r == '"':
			if inQuotes {
				inQuotes = false
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			} else {
				inQuotes = true
			}
		case unicode.IsSpace(r) && !inQuotes:
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
		default:
			currentToken.WriteRune(r)
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

func (r *REPL) printCommandUsage(commandName string) {
	command, exists := r.commands[commandName]
	if !exists {
		fmt.Println("Command not found.")
		return
	}

	fmt.Printf("%s ", command.Name)
	for _, arg := range command.Args {
		fmt.Printf("[%s:%s] ", arg.Name, arg.ArgType.Name())
	}
	fmt.Println("- ", command.Description)
}

// Help returns a handler function for the 'help' command that lists all available commands and their usage.
func (r *REPL) Help() func(args map[string]interface{}) error {
	return func(args map[string]interface{}) error {
		for name, cmd := range r.commands {
			fmt.Printf("%s: %s\n", name, cmd.Description)
			r.printCommandUsage(name)
		}
		return nil
	}
}
