package gin_handler

// Instruction	represents an arithmetic expression statement
// @Description	The instruction type can be either of the "calc" or "print".
// @Description	Type "calc" defines which arithmetic operation(op) (multiplication, addition, subtraction) to perform on two entities(left,right) and which variable(var) to save the result to.
// @Description	The entity(left/right) can be either an int64 literal or a variable name.
// @Description	Type "print" specifies the name of the variable(var) whose value needs to be output. In this case, there is no need to be fill in the remaining fields(op,left,right).
type Instruction struct {
	Type      string  `json:"type" binding:"required" example:"calc"`
	Operation *string `json:"op" example:"+"`
	Variable  string  `json:"var" binding:"required" example:"x"`
	Left      any     `json:"left" swaggertype:"integer" example:"2"`
	Right     any     `json:"right" swaggertype:"integer" example:"2"`
}

type Item struct {
	Var   string `json:"var" example:"x"`
	Value int    `json:"value" example:"4"`
}
