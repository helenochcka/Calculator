package core

import (
	"context"
	"errors"
	"sync"
)

type UseCase struct {
	service Service
}

func NewUseCase(s Service) UseCase {
	return UseCase{service: s}
}

func (uc *UseCase) Execute(ctx *context.Context, instructions []Instruction) ([]Item, error) {

	var items []Item
	var varsToPrint []string

	var wg sync.WaitGroup

	for _, instruction := range instructions {
		if instruction.Type == "calc" {
			if uc.service.GetResult(instruction.Result) != nil {
				return nil, errors.New("var already used")
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				result := uc.service.Algorithm(ctx, &instruction)
				if result != nil {
					uc.service.WriteResult(instruction.Result, *result)
					if uc.service.GetDependentValues(instruction.Result) != nil {
						wg.Add(1)
						go func() {
							defer wg.Done()
							uc.calcOfDependentVars(ctx, instruction.Result, &wg)
						}()

					}
				}
			}()

		} else if instruction.Type == "print" {
			varsToPrint = append(varsToPrint, instruction.Result)
		} else {
			return nil, errors.New("unknown type")
		}
	}

	wg.Wait()

	for _, varToPrint := range varsToPrint {
		items = append(items, Item{Var: varToPrint, Value: *uc.service.GetResult(varToPrint)})
	}

	return items, nil
}

func (uc *UseCase) calcOfDependentVars(ctx *context.Context, calcResult string, wg *sync.WaitGroup) {
	for _, instruction := range *uc.service.GetDependentValues(calcResult) {
		result := uc.service.Algorithm(ctx, &instruction)
		if result != nil {
			uc.service.WriteResult(instruction.Result, *result)
			if uc.service.GetDependentValues(instruction.Result) != nil {
				wg.Add(1)
				go func() {
					defer wg.Done()
					uc.calcOfDependentVars(ctx, instruction.Result, wg)
				}()
			}
		}
	}
	uc.service.DeleteKey(calcResult)
}
