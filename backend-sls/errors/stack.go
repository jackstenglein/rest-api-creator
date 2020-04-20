package errors

import "fmt"

// errorStack provides a convienent way to create and read string descriptions of
// an error's call stack.
type errorStack struct {
	// lines contains each line of the error stack
	lines []string
}

// newErrorStack returns an empty errorStack.
func newErrorStack() *errorStack {
	lines := make([]string, 0, 3)
	return &errorStack{lines}
}

// push adds line to the errorStack. If the errorStack is nil, push has no effect.
func (stack *errorStack) push(line string) {
	if stack != nil {
		stack.lines = append(stack.lines, line)
	}
}

// pop removes an error description from the errorStack. If the errorStack is nil or
// has no remaining descriptions, the empty string is returned.
func (stack *errorStack) pop() string {
	if stack == nil || len(stack.lines) == 0 {
		return ""
	}
	result := stack.lines[len(stack.lines)-1]
	stack.lines = stack.lines[:len(stack.lines)-1]
	return result
}

// hasElements returns true if the errorStack is not nil and has remaining error descriptions.
func (stack *errorStack) hasElements() bool {
	if stack == nil {
		return false
	}
	return len(stack.lines) > 0
}

// stack creates an errorStack that describes err. If err is nil, an empty errorStack is returned.
func stack(err error) *errorStack {
	stack := newErrorStack()
	cause := Cause(err)
	for err != cause && err != nil {
		line := fmt.Sprintf("%s: %s", Location(err), Message(err))
		stack.push(line)
		err = Previous(err)
	}
	stack.push(Message(err))
	return stack
}
