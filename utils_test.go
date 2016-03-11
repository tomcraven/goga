package ga_test

import (
	. "gopkg.in/check.v1"
)

// -----------------------------------------------------------------------
// IsTrue checker.

type isTrueChecker struct {
	*CheckerInfo
}

// The IsTrue checker tests whether the obtained value is true.
//
// For example:
//
//    c.Assert(val, IsTrue)
//
var IsTrue Checker = &isTrueChecker{
	&CheckerInfo{Name: "IsTrue", Params: []string{"value"}},
}

func (checker *isTrueChecker) Check(params []interface{}, names []string) (result bool, error string) {
	return isTrue(params[0]), ""
}

func isTrue(obtained interface{}) (result bool) {
	return (obtained == true)
}

// -----------------------------------------------------------------------
// IsFalse checker.

type isFalseChecker struct {
	*CheckerInfo
}

// The IsFalse checker tests whether the obtained value is true.
//
// For example:
//
//    c.Assert(val, IsFalse)
//
var IsFalse Checker = &isFalseChecker{
	&CheckerInfo{Name: "IsFalse", Params: []string{"value"}},
}

func (checker *isFalseChecker) Check(params []interface{}, names []string) (result bool, error string) {
	return isFalse(params[0]), ""
}

func isFalse(obtained interface{}) (result bool) {
	return (obtained == false)
}
