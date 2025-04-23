package service

import (
	"Calculator/api/arithmeticpb"
	"Calculator/internal/executor"
	"context"
	"strconv"
)

type InstructionStorage interface {
	Get(key string) *[]executor.Instruction
	Insert(key string, value executor.Instruction)
	Delete(key string)
	Clear()
}
type ResultStorage interface {
	Get(key string) *int
	Insert(key string, value int)
	Clear()
}
type Service struct {
	instructionStorage InstructionStorage
	resultStorage      ResultStorage
	client             arithmeticpb.ArithmeticServiceClient
}

func NewService(is InstructionStorage, rs ResultStorage, c arithmeticpb.ArithmeticServiceClient) Service {
	return Service{instructionStorage: is, resultStorage: rs, client: c}
}

func (s *Service) Algorithm(instruction *executor.Instruction) *int {
	left := s.castToInt(*instruction.Left)
	right := s.castToInt(*instruction.Right)

	if left == nil {
		if s.GetResult(*instruction.Left) == nil {
			if s.cyclicCheck(*instruction.Left, instruction.Result) {
				s.instructionStorage.Insert(*instruction.Left, *instruction)
			}
		} else {
			left = s.GetResult(*instruction.Left)
		}
	}
	if right == nil {
		if s.GetResult(*instruction.Right) == nil {
			if s.cyclicCheck(*instruction.Right, instruction.Result) {
				s.instructionStorage.Insert(*instruction.Right, *instruction)
			}
		} else {
			right = s.GetResult(*instruction.Right)
		}
	}
	if left != nil && right != nil {
		req := arithmeticpb.CalculationData{
			Op:    *instruction.Operation,
			Left:  int64(*left),
			Right: int64(*right),
		}

		response, err := s.client.Calculate(context.Background(), &req)

		if err == nil {
			result := int(response.Result)
			return &result
		}

	}
	return nil
}

func (s *Service) GetResult(key string) *int {
	return s.resultStorage.Get(key)
}

func (s *Service) WriteResult(key string, value int) {
	s.resultStorage.Insert(key, value)
}

func (s *Service) GetDependentValues(key string) *[]executor.Instruction {
	return s.instructionStorage.Get(key)
}

func (s *Service) DeleteKey(key string) {
	s.instructionStorage.Delete(key)
}

func (s *Service) ClearResults() {
	s.resultStorage.Clear()
}

func (s *Service) ClearDeps() {
	s.instructionStorage.Clear()
}

func (s *Service) cyclicCheck(dependentVar, calculatedVar string) bool {
	if s.GetDependentValues(dependentVar) != nil {
		for _, dependence := range *s.GetDependentValues(dependentVar) {
			if dependence.Result == calculatedVar {
				return false
			}
		}
	}
	return true
}

func (s *Service) castToInt(value string) *int {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &intValue
}
