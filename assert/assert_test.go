package assert

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/eugenetriguba/checkmate"
	"github.com/eugenetriguba/checkmate/internal/cmtest"
)

type assertFn func(t checkmate.TestingT, args []any)

func wrappedAssertNil(t checkmate.TestingT, args []any) {
	if len(args) > 1 {
		Nil(t, args[0], args[1:]...)
	} else {
		Nil(t, args[0])
	}
}

func wrappedAssertNotNil(t checkmate.TestingT, args []any) {
	if len(args) > 1 {
		NotNil(t, args[0], args[1:]...)
	} else {
		NotNil(t, args[0])
	}
}

func wrappedAssertTrue(t checkmate.TestingT, args []any) {
	if len(args) > 1 {
		True(t, args[0].(bool), args[1:]...)
	} else {
		True(t, args[0].(bool))
	}
}

func wrappedAssertFalse(t checkmate.TestingT, args []any) {
	if len(args) > 1 {
		False(t, args[0].(bool), args[1:]...)
	} else {
		False(t, args[0].(bool))
	}
}

func wrappedAssertDeepEqual(t checkmate.TestingT, args []any) {
	if len(args) > 2 {
		DeepEqual(t, args[0], args[1], args[2:]...)
	} else {
		DeepEqual(t, args[0], args[1])
	}
}

func wrappedAssertNotDeepEqual(t checkmate.TestingT, args []any) {
	if len(args) > 2 {
		NotDeepEqual(t, args[0], args[1], args[2:]...)
	} else {
		NotDeepEqual(t, args[0], args[1])
	}
}

func wrappedAssertErrorIs(t checkmate.TestingT, args []any) {
	if len(args) > 2 {
		ErrorIs(t, args[0].(error), args[1].(error), args[2:]...)

	} else {
		ErrorIs(t, args[0].(error), args[1].(error))
	}
}

func wrappedAssertNotErrorIs(t checkmate.TestingT, args []any) {
	if len(args) > 2 {
		NotErrorIs(t, args[0].(error), args[1].(error), args[2:]...)
	} else {
		NotErrorIs(t, args[0].(error), args[1].(error))
	}
}

func wrappedAssertErrorContains(t checkmate.TestingT, args []any) {
	if len(args) > 2 {
		ErrorContains(t, args[0].(error), args[1].(string), args[2:]...)
	} else {
		ErrorContains(t, args[0].(error), args[1].(string))
	}
}

func wrappedAssertNotErrorContains(t checkmate.TestingT, args []any) {
	if len(args) > 2 {
		NotErrorContains(t, args[0].(error), args[1].(string), args[2:]...)
	} else {
		NotErrorContains(t, args[0].(error), args[1].(string))
	}
}

func wrappedAssertEqual(t checkmate.TestingT, args []any) {
	if len(args) > 2 {
		Equal(t, args[0], args[1], args[2:]...)
	} else {
		Equal(t, args[0], args[1])
	}
}

func wrappedAssertNotEqual(t checkmate.TestingT, args []any) {
	if len(args) > 2 {
		NotEqual(t, args[0], args[1], args[2:]...)
	} else {
		NotEqual(t, args[0], args[1])
	}
}

var passingTestFns = []struct {
	name        string
	assertionFn assertFn
	args        []any
}{
	{"AssertNil", wrappedAssertNil, []any{nil}},
	{"AssertNotNil", wrappedAssertNotNil, []any{1}},
	{"AssertTrue", wrappedAssertTrue, []any{true}},
	{"AssertFalse", wrappedAssertFalse, []any{false}},
	{"AssertErrorIs", wrappedAssertErrorIs, []any{os.ErrExist, os.ErrExist}},
	{"AssertNotErrorIs", wrappedAssertNotErrorIs, []any{os.ErrExist, os.ErrClosed}},
	{"AssertErrorContains", wrappedAssertErrorContains, []any{errors.New("error 1"), "error 1"}},
	{"AssertNotErrorContains", wrappedAssertNotErrorContains, []any{errors.New("error 1"), "error 2"}},
	{"AssertDeepEqual", wrappedAssertDeepEqual, []any{5, 5}},
	{"AssertNotDeepEqual", wrappedAssertNotDeepEqual, []any{5, 6}},
	{"AssertEqual", wrappedAssertEqual, []any{5, 5}},
	{"AssertNotEqual", wrappedAssertNotEqual, []any{5, 10}},
}

var failingTestFns = []struct {
	name        string
	assertionFn assertFn
	args        []any
}{
	{"AssertNil", wrappedAssertNil, []any{1}},
	{"AssertNotNil", wrappedAssertNotNil, []any{nil}},
	{"AssertTrue", wrappedAssertTrue, []any{false}},
	{"AssertFalse", wrappedAssertFalse, []any{true}},
	{"AssertErrorIs", wrappedAssertErrorIs, []any{errors.New("error 1"), errors.New("error 2")}},
	{"AssertNotErrorIs", wrappedAssertNotErrorIs, []any{os.ErrClosed, os.ErrClosed}},
	{"AssertErrorContains", wrappedAssertErrorContains, []any{errors.New("error 1"), "error 2"}},
	{"AssertNotErrorContains", wrappedAssertNotErrorContains, []any{errors.New("error 1"), "error 1"}},
	{"AssertDeepEqual", wrappedAssertDeepEqual, []any{5, 10}},
	{"AssertNotDeepEqual", wrappedAssertNotDeepEqual, []any{5, 5}},
	{"AssertEqual", wrappedAssertEqual, []any{5, 10}},
	{"AssertNotEqual", wrappedAssertNotEqual, []any{5, 5}},
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

				testFn.assertionFn(mockT, append(testFn.args, testCase.args...))

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

func TestAssertFnsCallHelper(t *testing.T) {
	for _, testFn := range passingTestFns {
		t.Run(testFn.name, func(t *testing.T) {
			mockHelperT := &cmtest.MockHelperT{}

			testFn.assertionFn(mockHelperT, testFn.args)

			if !mockHelperT.HelperCalled {
				t.Errorf("%s: expected HelperT to be called", testFn.name)
			}
		})
	}
}

func TestAssertFnsShouldNotFailOrEmitLogsOnSuccess(t *testing.T) {
	for _, testFn := range passingTestFns {
		t.Run(testFn.name, func(t *testing.T) {
			mockT := &cmtest.MockT{}
			testFn.assertionFn(mockT, testFn.args)

			if mockT.FailNowCalled {
				t.Fatalf(
					"%s: fail now called, expected the test to not be failed on true condition",
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
		{"UnequalIntegers", 5, 10, true, []string{"expected 5 to equal 10"}},
		{"EqualFloats", 5.123, 5.123, false, []string{}},
		{"UnequalFloats", 5.123, 5.1234, true, []string{"expected 5.123 to equal 5.1234"}},
		{"EqualStrings", "test", "test", false, []string{}},
		{"UnequalStrings", "test", "fail", true, []string{"expected test to equal fail"}},
		{"EqualBooleans", true, true, false, []string{}},
		{"UnequalBooleans", false, true, true, []string{"expected false to equal true"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockT := &cmtest.MockT{}
			Equal(mockT, tc.actual, tc.expected)

			if mockT.FailNowCalled != tc.shouldFail {
				t.Errorf(
					"%s: FailNowCalled = %v, want %v",
					tc.name, mockT.FailNowCalled, tc.shouldFail,
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
		if mockT.FailNowCalled {
			t.Error("AssertDeepEqual failed when it should have passed")
		}
	})

	t.Run("Not equal structs", func(t *testing.T) {
		mockT := &cmtest.MockT{}
		DeepEqual(mockT, TestStruct{"Alice", 30}, TestStruct{"Bob", 30})
		if !mockT.FailNowCalled {
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
