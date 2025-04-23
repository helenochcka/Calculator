package use_case

import (
	"Calculator/internal/executor"
	"Calculator/internal/executor/service"
	"errors"
	"sync"
)

type UseCase struct {
	service service.Service
}

func NewUseCase(s service.Service) UseCase {
	return UseCase{service: s}
}

func (uc *UseCase) Execute(instructions []executor.Instruction) ([]executor.Item, error) {

	var items []executor.Item
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
				result := uc.service.Algorithm(&instruction)
				if result != nil {
					uc.service.WriteResult(instruction.Result, *result)
					if uc.service.GetDependentValues(instruction.Result) != nil {
						wg.Add(1)
						go func() {
							defer wg.Done()
							uc.calcOfDependentVars(instruction.Result, &wg)
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
		items = append(items, executor.Item{Var: varToPrint, Value: *uc.service.GetResult(varToPrint)})
	}

	uc.service.ClearResults()
	uc.service.ClearDeps()

	return items, nil
}

func (uc *UseCase) calcOfDependentVars(calcResult string, wg *sync.WaitGroup) {
	for _, instruction := range *uc.service.GetDependentValues(calcResult) {
		result := uc.service.Algorithm(&instruction)
		if result != nil {
			uc.service.WriteResult(instruction.Result, *result)
			if uc.service.GetDependentValues(instruction.Result) != nil {
				wg.Add(1)
				go func() {
					defer wg.Done()
					uc.calcOfDependentVars(instruction.Result, wg)
				}()
			}
		}
	}
	uc.service.DeleteKey(calcResult)
}
