package checkmate

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

// AssertDeepEqual checks if two values are deeply equal using Google's
// go-cmp cmp.Equal. If they are not equal, it logs the differences
// using cmp.Diff. An optional message and arguments for any format
// placeholders in that message can be provided for if a failure occurs.
func AssertDeepEqual(t TestingT, actual, expected any, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		if len(msgAndArgs) == 0 {
			msgAndArgs = []any{"mismatch (-expected +actual):\n%s", diff}
		}
		Assert(t, false, msgAndArgs...)
	}
}

// AssertEqual checks if two primitive values are equal. An optional
// message and arguments for any format placeholders in that message
// can be provided for if a failure occurs.
func AssertEqual(t TestingT, actual, expected any, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if actual != expected {
		if len(msgAndArgs) == 0 {
			msgAndArgs = []any{"expected %v to equal %v", expected, actual}
		}
		Assert(t, false, msgAndArgs...)
	}
}

// AssertNotEqual checks if two values are not equal. It fails the test if
// the values are equal. An optional message and arguments for any format
// placeholders in that message can be provided if a failure occurs.
func AssertNotEqual(t TestingT, actual, expected any, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if actual == expected {
		if len(msgAndArgs) == 0 {
			msgAndArgs = []any{"expected %v to not equal %v", expected, actual}
		}
		Assert(t, false, msgAndArgs...)
	}
}

// AssertLenEqual checks if the length of a given slice or array equals the
// expected length. It fails the test if the lengths are not equal. An optional
// message and arguments for any format placeholders in that message can be provided
// if a failure occurs.
func AssertLenEqual[T any](t TestingT, l []T, expectedLen int, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"expected %v to have len %d, got len %d", l, len(l), expectedLen}
	}
	Assert(t, len(l) == expectedLen, msgAndArgs...)
}

// Assert checks a boolean condition and fails the test if the condition is false.
// This is a general-purpose assertion function used to validate if a condition holds true.
// An optional message and arguments for any format placeholders in that message
// can be provided if a failure occurs.
func Assert(t TestingT, condition bool, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}
	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"assertion failed"}
	}
	passed := Check(t, condition, msgAndArgs...)
	if !passed {
		t.FailNow()
	}
}

// Check evaluates a boolean condition and logs a message if the condition is false.
// Unlike Assert, it does not stop the test when the condition is false but returns
// a boolean indicating the result. This function is useful for cases where you want
// to perform multiple checks in a single test function without stopping the test
// execution after the first failure. An optional message and arguments for any format
// placeholders in that message can be provided if the check fails.
func Check(t TestingT, condition bool, msgAndArgs ...any) bool {
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
				t.Log("checkmate: assertion called with a non-string message, using default message")
			}
		}
		t.Log(message)
		return false
	}
	return true
}

// The subset of testing.T which is used by the
// checkmate package.
type TestingT interface {
	Log(args ...any)
	FailNow()
}

type helperT interface {
	Helper()
}
