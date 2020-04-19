package errors

import "fmt"

type ErrorStack struct {
	lines []string
}

func NewErrorStack() *ErrorStack {
	lines := make([]string, 0, 3)
	return &ErrorStack{lines}
}

func (stack *ErrorStack) Push(line string) {
	if stack != nil {
		stack.lines = append(stack.lines, line)
	}
}

func (stack *ErrorStack) Pop() string {
	if stack == nil || len(stack.lines) == 0 {
		return ""
	}
	result := stack.lines[len(stack.lines)-1]
	stack.lines = stack.lines[:len(stack.lines)-1]
	return result
}

func (stack *ErrorStack) HasElements() bool {
	if stack == nil {
		return false
	}
	return len(stack.lines) > 0
}

func Stack(err error) *ErrorStack {
	stack := NewErrorStack()
	cause := Cause(err)
	for err != cause && err != nil {
		line := fmt.Sprintf("%s: %s", Location(err), Message(err))
		stack.Push(line)
		err = Previous(err)
	}
	stack.Push(Message(err))
	return stack
}
