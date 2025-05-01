package executor

import "errors"

var ErrDeclaringQueue = errors.New("error when declaring queue")
var ErrCyclicDependency = errors.New("cyclic dependency detected")
var ErrConsumingResult = errors.New("error when consuming result")
var ErrUnknownInstructionType = errors.New("unknown type of instruction")
var ErrVarAlreadyUsed = errors.New("variable is already used")
var ErrCalcExpression = errors.New("error when calculating expression")
var ErrVarToPrintNotFound = errors.New("variable to print not found")
var ErrVarNeverBeCalc = errors.New("expression argument variable will never be calculated")
