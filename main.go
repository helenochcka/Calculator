package main

import (
	"Calculator/core"
	"fmt"
)

func main() {
	service := core.NewService()
	useCase := core.NewUseCase(service)

	op := "+"
	left := "1"
	right := "2"
	//res := "x"
	res2 := "y"

	dto1 := core.DTO{
		Type:      "calc",
		Operation: &op,
		Result:    "x",
		Left:      &res2,
		Right:     &right,
	}

	dto2 := core.DTO{
		Type:      "print",
		Operation: nil,
		Result:    "x",
		Left:      nil,
		Right:     nil,
	}

	//op3 := "+"
	dto3 := core.DTO{
		Type:      "calc",
		Operation: &op,
		Result:    "y",
		Left:      &left,
		Right:     &right,
	}

	dto4 := core.DTO{
		Type:      "print",
		Operation: nil,
		Result:    "y",
		Left:      nil,
		Right:     nil,
	}

	dtos := []core.DTO{dto1, dto2, dto3, dto4}

	items, err := useCase.Execute(dtos)

	if err != nil {
		fmt.Printf("Error %s", err)
	} else {
		println(items[0].Var, items[0].Value, items[1].Var, items[1].Value)
	}

}
