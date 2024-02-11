package check

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/eugenetriguba/checkmate"
	"github.com/google/go-cmp/cmp"
)

type helperT interface {
	Helper()
}

// Nil checks whether the value equals nil.
func Nil(t checkmate.TestingT, value any, msgAndArgs ...any) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"expected value to be nil, got %v", value}
	}

	isNil := false
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			isNil = true
		}
	} else if value == nil {
		isNil = true
	}

	return check(t, isNil, msgAndArgs...)
}

// NotNil checks whether the value does not equal nil.
func NotNil(t checkmate.TestingT, value any, msgAndArgs ...any) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"expected value to not be nil, got nil", value}
	}

	return check(t, value != nil, msgAndArgs...)
}

// True checks whether the condition is true.
func True(t checkmate.TestingT, condition bool, msgAndArgs ...any) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"expected condition to be true, got false"}
	}

	return check(t, condition, msgAndArgs...)
}

// False checks whether the condition is false.
func False(t checkmate.TestingT, condition bool, msgAndArgs ...any) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"expected condition to be false, got true"}
	}

	return check(t, !condition, msgAndArgs...)
}

// ErrorIs checks whether the target error occurs within err's error tree.
func ErrorIs(t checkmate.TestingT, err error, target error, msgAndArgs ...any) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"expected error %v to have error %v in its tree", err, target}
	}

	return check(t, errors.Is(err, target), msgAndArgs...)
}

// NotErrorIs checks whether the target error does not occur within
// the err's error tree.
func NotErrorIs(t checkmate.TestingT, err error, target error, msgAndArgs ...any) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"expected error %v to not have error %v in its tree", err, target}
	}

	return check(t, !errors.Is(err, target), msgAndArgs...)
}

// ErrorContains checks whether the given err contains the errText
// in the err.Error() output.
func ErrorContains(t checkmate.TestingT, err error, errText string, msgAndArgs ...any) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"expected err to contain %s, got %s", errText, err.Error()}
	}

	return check(t, strings.Contains(err.Error(), errText), msgAndArgs...)
}

// NotErrorContains checks whether the given err does not contain the errText
// in the err.Error() output.
func NotErrorContains(t checkmate.TestingT, err error, errText string, msgAndArgs ...any) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"expected err to contain not %s, got that it does", errText}
	}

	return check(t, !strings.Contains(err.Error(), errText), msgAndArgs...)
}

// DeepEqual checks if two values are deeply equal. If they are not equal,
// it logs the differences.
func DeepEqual(t checkmate.TestingT, actual, expected any, msgAndArgs ...any) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	diff := cmp.Diff(expected, actual)
	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"mismatch (-expected +actual):\n%s", diff}
	}

	return check(t, diff == "", msgAndArgs...)
}

// NotDeepEqual checks if two values are not deeply equal.
func NotDeepEqual(t checkmate.TestingT, actual, expected any, msgAndArgs ...any) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	diff := cmp.Diff(expected, actual)
	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"expected %v to not equal %v, got that they're equal", actual, expected}
	}
	return check(t, diff != "", msgAndArgs...)
}

// Equal checks if two primitive values are equal.
func Equal(t checkmate.TestingT, actual, expected any, msgAndArgs ...any) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"expected %v to equal %v", expected, actual}
	}

	return check(t, actual == expected, msgAndArgs...)
}

// NotEqual checks if two values are not equal. It fails the test if
// the values are equal.
func NotEqual(t checkmate.TestingT, actual, expected any, msgAndArgs ...any) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"expected %v to not equal %v", actual, expected}
	}

	return check(t, actual != expected, msgAndArgs...)
}

// Check evaluates a boolean condition and if the condition is false,
// it will log out a message and mark the test as failed. However, it does
// not immediately stop execution, unlike the assert functions.
func check(t checkmate.TestingT, condition bool, msgAndArgs ...any) bool {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if !condition {
		message := "check failed"
		if len(msgAndArgs) > 0 {
			switch msgAndArgs[0].(type) {
			case string:
				message = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
			default:
				t.Log("checkmate: check called with a non-string message, using default message")
			}
		}
		t.Log(message)
		t.Fail()
	}

	return condition
}
