package executor

import "errors"

var ErrDeclaringQueue = errors.New("error when declaring queue")
var ErrCyclicDependency = errors.New("cyclic dependency detected")
var ErrConsumingResult = errors.New("error when consuming result")
var ErrVarAlreadyUsed = errors.New("variable is already used")
var ErrCalcExpression = errors.New("error when calculating expression")
var ErrVarWillNeverBeCalc = errors.New("variable will never be calculated")
var ErrReqIdMissing = errors.New("request id is missing in the context")
var ErrUnsupportedArgType = errors.New("unsupported argument type")
