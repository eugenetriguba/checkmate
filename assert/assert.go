package assert

import (
	"github.com/eugenetriguba/checkmate"
	"github.com/eugenetriguba/checkmate/check"
)

// Nil asserts whether the value equals nil.
func Nil(t checkmate.TestingT, value any, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if passed := check.Nil(t, value, msgAndArgs...); !passed {
		t.FailNow()
	}
}

// NotNil asserts whether the value does not equal nil.
func NotNil(t checkmate.TestingT, value any, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if passed := check.NotNil(t, value, msgAndArgs...); !passed {
		t.FailNow()
	}
}

// True asserts whether the condition is true.
func True(t checkmate.TestingT, condition bool, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if passed := check.True(t, condition, msgAndArgs...); !passed {
		t.FailNow()
	}
}

// False asserts whether the condition is false.
func False(t checkmate.TestingT, condition bool, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if passed := check.False(t, condition, msgAndArgs...); !passed {
		t.FailNow()
	}
}

// ErrorIs asserts whether the target error occurs within err's error tree.
func ErrorIs(t checkmate.TestingT, err error, target error, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if passed := check.ErrorIs(t, err, target, msgAndArgs...); !passed {
		t.FailNow()
	}
}

// NotErrorIs asserts whether the target error does not occur within
// the err's error tree.
func NotErrorIs(t checkmate.TestingT, err error, target error, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if passed := check.NotErrorIs(t, err, target, msgAndArgs...); !passed {
		t.FailNow()
	}
}

// ErrorContains asserts whether the given err contains the errText in the err.Error() output.
func ErrorContains(t checkmate.TestingT, err error, errText string, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if passed := check.ErrorContains(t, err, errText, msgAndArgs...); !passed {
		t.FailNow()
	}
}

// NotErrorContains asserts whether the given err does not contain the errText
// in the err.Error() output.
func NotErrorContains(t checkmate.TestingT, err error, errText string, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if passed := check.NotErrorContains(t, err, errText, msgAndArgs...); !passed {
		t.FailNow()
	}
}

// DeepEqual asserts that two values are deeply equal. If they are not equal,
// it logs the differences.
func DeepEqual(t checkmate.TestingT, actual, expected any, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if passed := check.DeepEqual(t, actual, expected, msgAndArgs...); !passed {
		t.FailNow()
	}
}

// NotDeepEqual asserts if two values are not deeply equal.
func NotDeepEqual(t checkmate.TestingT, actual, expected any, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if passed := check.NotDeepEqual(t, actual, expected, msgAndArgs...); !passed {
		t.FailNow()
	}
}

// Equal asserts if two primitive values are equal.
func Equal(t checkmate.TestingT, actual, expected any, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if passed := check.Equal(t, actual, expected, msgAndArgs...); !passed {
		t.FailNow()
	}
}

// NotEqual asserts if two values are not equal.
func NotEqual(t checkmate.TestingT, actual, expected any, msgAndArgs ...any) {
	if ht, ok := t.(helperT); ok {
		ht.Helper()
	}

	if passed := check.NotEqual(t, actual, expected, msgAndArgs...); !passed {
		t.FailNow()
	}
}

type helperT interface {
	Helper()
}
