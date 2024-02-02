package check

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/eugenetriguba/checkmate"
	"github.com/eugenetriguba/checkmate/internal/cmtest"
)

type checkFn func(t checkmate.TestingT, args []any) bool

func wrappedCheckNil(t checkmate.TestingT, args []any) bool {
	if len(args) > 1 {
		return Nil(t, args[0], args[1:]...)
	} else {
		return Nil(t, args[0])
	}
}

func wrappedCheckNotNil(t checkmate.TestingT, args []any) bool {
	if len(args) > 1 {
		return NotNil(t, args[0], args[1:]...)
	} else {
		return NotNil(t, args[0])
	}
}

func wrappedCheckTrue(t checkmate.TestingT, args []any) bool {
	if len(args) > 1 {
		return True(t, args[0].(bool), args[1:]...)
	} else {
		return True(t, args[0].(bool))
	}
}

func wrappedCheckFalse(t checkmate.TestingT, args []any) bool {
	if len(args) > 1 {
		return False(t, args[0].(bool), args[1:]...)
	} else {
		return False(t, args[0].(bool))
	}
}

func wrappedCheckDeepEqual(t checkmate.TestingT, args []any) bool {
	if len(args) > 2 {
		return DeepEqual(t, args[0], args[1], args[2:]...)
	} else {
		return DeepEqual(t, args[0], args[1])
	}
}

func wrappedCheckNotDeepEqual(t checkmate.TestingT, args []any) bool {
	if len(args) > 2 {
		return NotDeepEqual(t, args[0], args[1], args[2:]...)
	} else {
		return NotDeepEqual(t, args[0], args[1])
	}
}

func wrappedCheckErrorIs(t checkmate.TestingT, args []any) bool {
	if len(args) > 2 {
		return ErrorIs(t, args[0].(error), args[1].(error), args[2:]...)
	} else {
		return ErrorIs(t, args[0].(error), args[1].(error))
	}
}

func wrappedCheckNotErrorIs(t checkmate.TestingT, args []any) bool {
	if len(args) > 2 {
		return NotErrorIs(t, args[0].(error), args[1].(error), args[2:]...)
	} else {
		return NotErrorIs(t, args[0].(error), args[1].(error))
	}
}

func wrappedCheckErrorContains(t checkmate.TestingT, args []any) bool {
	if len(args) > 2 {
		return ErrorContains(t, args[0].(error), args[1].(string), args[2:]...)
	} else {
		return ErrorContains(t, args[0].(error), args[1].(string))
	}
}

func wrappedCheckNotErrorContains(t checkmate.TestingT, args []any) bool {
	if len(args) > 2 {
		return NotErrorContains(t, args[0].(error), args[1].(string), args[2:]...)
	} else {
		return NotErrorContains(t, args[0].(error), args[1].(string))
	}
}

func wrappedCheckEqual(t checkmate.TestingT, args []any) bool {
	if len(args) > 2 {
		return Equal(t, args[0], args[1], args[2:]...)
	} else {
		return Equal(t, args[0], args[1])
	}
}

func wrappedCheckNotEqual(t checkmate.TestingT, args []any) bool {
	if len(args) > 2 {
		return NotEqual(t, args[0], args[1], args[2:]...)
	} else {
		return NotEqual(t, args[0], args[1])
	}
}

var passingTestFns = []struct {
	name string
	fn   checkFn
	args []any
}{
	{"CheckNil", wrappedCheckNil, []any{nil}},
	{"CheckNotNil", wrappedCheckNotNil, []any{1}},
	{"CheckTrue", wrappedCheckTrue, []any{true}},
	{"CheckFalse", wrappedCheckFalse, []any{false}},
	{"CheckErrorIs", wrappedCheckErrorIs, []any{os.ErrExist, os.ErrExist}},
	{"CheckNotErrorIs", wrappedCheckNotErrorIs, []any{os.ErrExist, os.ErrClosed}},
	{"CheckErrorContains", wrappedCheckErrorContains, []any{errors.New("error 1"), "error 1"}},
	{"CheckNotErrorContains", wrappedCheckNotErrorContains, []any{errors.New("error 1"), "error 2"}},
	{"CheckDeepEqual", wrappedCheckDeepEqual, []any{5, 5}},
	{"CheckNotDeepEqual", wrappedCheckNotDeepEqual, []any{5, 6}},
	{"CheckEqual", wrappedCheckEqual, []any{5, 5}},
	{"CheckNotEqual", wrappedCheckNotEqual, []any{5, 10}},
}

var failingTestFns = []struct {
	name string
	fn   checkFn
	args []any
}{
	{"CheckNil", wrappedCheckNil, []any{1}},
	{"CheckNotNil", wrappedCheckNotNil, []any{nil}},
	{"CheckTrue", wrappedCheckTrue, []any{false}},
	{"CheckFalse", wrappedCheckFalse, []any{true}},
	{"CheckErrorIs", wrappedCheckErrorIs, []any{errors.New("error 1"), errors.New("error 2")}},
	{"CheckNotErrorIs", wrappedCheckNotErrorIs, []any{os.ErrClosed, os.ErrClosed}},
	{"CheckErrorContains", wrappedCheckErrorContains, []any{errors.New("error 1"), "error 2"}},
	{"CheckNotErrorContains", wrappedCheckNotErrorContains, []any{errors.New("error 1"), "error 1"}},
	{"CheckDeepEqual", wrappedCheckDeepEqual, []any{5, 10}},
	{"CheckNotDeepEqual", wrappedCheckNotDeepEqual, []any{5, 5}},
	{"CheckEqual", wrappedCheckEqual, []any{5, 10}},
	{"CheckNotEqual", wrappedCheckNotEqual, []any{5, 5}},
}

func TestOptionalMessageAndArgs(t *testing.T) {
	testCases := []struct {
		name         string
		args         []any
		expectedLogs []string
	}{
		{"Plain message", []any{"my message"}, []string{"my message"}},
		{"Message with format placeholders", []any{"my message: %d", 5}, []string{"my message: 5"}},
		{
			"Invalid message type",
			[]any{5},
			[]string{
				"checkmate: check called with a non-string message, using default message",
				"check failed",
			},
		},
	}

	for _, testFn := range failingTestFns {
		for _, testCase := range testCases {
			testName := fmt.Sprintf("%s: %s", testFn.name, testCase.name)
			t.Run(testName, func(t *testing.T) {
				mockT := &cmtest.MockT{}

				testFn.fn(mockT, append(testFn.args, testCase.args...))

				if len(mockT.Logs) != len(testCase.expectedLogs) {
					t.Errorf(
						"%s: expected %d (%v) log messages but got %d (%v)",
						testName, len(testCase.expectedLogs),
						testCase.expectedLogs, len(mockT.Logs), mockT.Logs,
					)
				}
				if len(mockT.Logs) == len(testCase.expectedLogs) {
					for i, expectedLogMsg := range testCase.expectedLogs {
						if mockT.Logs[i] != expectedLogMsg {
							t.Errorf(
								"%s: expected log message '%s', got '%s'",
								testName, expectedLogMsg, mockT.Logs[i],
							)
						}
					}
				}
			})
		}
	}
}

func TestCheckFnsCallHelper(t *testing.T) {
	for _, testFn := range passingTestFns {
		t.Run(testFn.name, func(t *testing.T) {
			mockHelperT := &cmtest.MockHelperT{}

			testFn.fn(mockHelperT, testFn.args)

			if !mockHelperT.HelperCalled {
				t.Errorf("%s: expected HelperT to be called", testFn.name)
			}
		})
	}
}

func TestCheckFnsShouldNotFailOrEmitLogsOnSuccess(t *testing.T) {
	for _, testFn := range passingTestFns {
		t.Run(testFn.name, func(t *testing.T) {
			mockT := &cmtest.MockT{}
			testFn.fn(mockT, testFn.args)

			if mockT.FailNowCalled || mockT.FailCalled {
				t.Fatalf(
					"%s: fail called, expected the test to not be failed on true condition",
					testFn.name,
				)
			}
			if len(mockT.Logs) > 0 {
				t.Fatalf(
					"%s: %v logs emitted, expected the test to not emit any logs on true condition",
					testFn.name, mockT.Logs,
				)
			}
		})
	}
}

func TestAssertEqual(t *testing.T) {
	testCases := []struct {
		name        string
		actual      any
		expected    any
		shouldFail  bool
		logMessages []string
	}{
		{"EqualIntegers", 5, 5, false, []string{}},
		{"UnequalIntegers", 5, 10, true, []string{"expected 10 to equal 5"}},
		{"EqualFloats", 5.123, 5.123, false, []string{}},
		{"UnequalFloats", 5.123, 5.1234, true, []string{"expected 5.1234 to equal 5.123"}},
		{"EqualStrings", "test", "test", false, []string{}},
		{"UnequalStrings", "test", "fail", true, []string{"expected fail to equal test"}},
		{"EqualBooleans", true, true, false, []string{}},
		{"UnequalBooleans", false, true, true, []string{"expected true to equal false"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockT := &cmtest.MockT{}
			Equal(mockT, tc.actual, tc.expected)

			if mockT.FailCalled != tc.shouldFail {
				t.Errorf(
					"%s: FailCalled = %v, want %v",
					tc.name, mockT.FailCalled, tc.shouldFail,
				)
			}
			if len(mockT.Logs) != len(tc.logMessages) {
				t.Errorf(
					"%s: expected %d (%v) log messages but got %d (%v)",
					tc.name, len(mockT.Logs), mockT.Logs, len(tc.logMessages), tc.logMessages,
				)
			}
			if len(mockT.Logs) == len(tc.logMessages) {
				for i, expectedLogMsg := range tc.logMessages {
					if mockT.Logs[i] != expectedLogMsg {
						t.Errorf(
							"%s: expected log message '%s', got '%s'",
							tc.name, expectedLogMsg, mockT.Logs[0],
						)
					}
				}
			}
		})
	}
}

func TestAssertDeepEqual(t *testing.T) {
	type TestStruct struct {
		Name string
		Age  int
	}

	t.Run("Equal structs", func(t *testing.T) {
		mockT := &cmtest.MockT{}
		DeepEqual(mockT, TestStruct{"Alice", 30}, TestStruct{"Alice", 30})
		if mockT.FailCalled {
			t.Error("AssertDeepEqual failed when it should have passed")
		}
	})

	t.Run("Not equal structs", func(t *testing.T) {
		mockT := &cmtest.MockT{}
		DeepEqual(mockT, TestStruct{"Alice", 30}, TestStruct{"Bob", 30})
		if !mockT.FailCalled {
			t.Error("AssertDeepEqual passed when it should have failed")
		}
		if len(mockT.Logs) == 0 || !containsDiffMessage(mockT.Logs[0]) {
			t.Error("AssertDeepEqual did not log the correct diff message")
		}
	})
}

func containsDiffMessage(log string) bool {
	return strings.Contains(log, "(-expected +actual)")
}

func TestAssertErrorIsChecksErrTree(t *testing.T) {
	mockT := &cmtest.MockT{}
	actual := fmt.Errorf("wrapped: %w", os.ErrInvalid)
	expected := os.ErrInvalid

	ErrorIs(mockT, actual, expected)

	if mockT.FailNowCalled {
		t.Fatal("AssertErrorIs failed when it should have passed")
	}
}

func TestCheckReturnValues(t *testing.T) {
	t.Run("should return true on true condition", func(t *testing.T) {
		mockT := &cmtest.MockT{}

		passed := check(mockT, true)

		if !passed {
			t.Fatal("Check should have returned true on true condition")
		}
	})

	t.Run("should return false on false condition", func(t *testing.T) {
		mockT := &cmtest.MockT{}

		passed := check(mockT, false)

		if passed {
			t.Fatal("Check should have returned false on false condition")
		}
	})
}
