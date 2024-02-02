// Checkmate is a library designed to enhance the testing experience by providing a
// set of assertion functions. With Go, there is no built-in assertion library which
// can make writing tests feel more verbose than it needs to be, so Checkmate aims to
// fill that gap.
//
// This library is divided into two different packages. It provides a `check` package and an
// `assert` package. The difference between them is that the functions inside of the `assert`
// package will immediately fail the test on an assertion failure. However, the functions inside
// of the `check` package will mark the test as failed and return a boolean result, however, they
// will continue execution until the test has finished. Outside of those two differences, the
// API and behavior is the same between them.
//
// Furthermore, any of the assertion of check package will accept a variadic argument at the
// end called `msgAndArgs`. This allows the caller to pass in their own custom message on test
// failure along with any arguments for the message if any format placeholders were used.
package checkmate

// The subset of testing.T which is used by the
// checkmate package.
type TestingT interface {
	Log(args ...any)
	Fail()
	FailNow()
}
