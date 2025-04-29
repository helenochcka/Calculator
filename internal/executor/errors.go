package executor

import "errors"

var ErrDeclaringQueue = errors.New("error declaring queue")
var CyclicDependency = errors.New("cyclic dependency detected")
var ErrConsumingResult = errors.New("error consuming result")
var UnknownTypeOfInstruction = errors.New("unknown type of instruction")
var VarIsAlreadyUsed = errors.New("variable is already used")
var ErrCalcExpression = errors.New("error calculating expression")
var VarToPrintNotFound = errors.New("variable to print not found")
var VarNeverBeCalc = errors.New("variable on which other variables depend will never be calculated")
