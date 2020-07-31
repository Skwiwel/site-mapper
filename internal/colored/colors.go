package colored

import "fmt"

func color(colorString string) func(...interface{}) string {
	return func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
}

var (
	Warning = color("\033[33;5m%s\033[0m")
	Grey    = color("\033[33;90m%s\033[0m")
)
