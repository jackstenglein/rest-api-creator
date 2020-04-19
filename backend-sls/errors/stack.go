package errors

import "fmt"

type errorStack struct {
	lines []string
}

func newErrorStack() *errorStack {
	lines := make([]string, 0, 3)
	return &errorStack{lines}
}

func (stack *errorStack) push(line string) {
	if stack != nil {
		stack.lines = append(stack.lines, line)
	}
}

func (stack *errorStack) pop() string {
	if stack == nil || len(stack.lines) == 0 {
		return ""
	}
	result := stack.lines[len(stack.lines)-1]
	stack.lines = stack.lines[:len(stack.lines)-1]
	return result
}

func (stack *errorStack) hasElements() bool {
	if stack == nil {
		return false
	}
	return len(stack.lines) > 0
}

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
