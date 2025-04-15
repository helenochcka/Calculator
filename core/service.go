package core

import (
	gen "Calculator/another_service/grpc/gen"
	"context"
	"fmt"
	"strconv"
	//"time"
)

type InstructionStorage interface {
	Get(key string) *[]Instruction
	Insert(key string, value Instruction)
	Delete(key string)
}
type ResultStorage interface {
	Get(key string) *int
	Insert(key string, value int)
}
type Service struct {
	instructionStorage InstructionStorage
	resultStorage      ResultStorage
	client             gen.CalcServiceClient
}

func NewService(is InstructionStorage, rs ResultStorage, c gen.CalcServiceClient) Service {
	return Service{instructionStorage: is, resultStorage: rs, client: c}
}

func (s *Service) Algorithm(ctx *context.Context, instruction *Instruction) *int {
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
		req := gen.CalcRequest{
			Op:    *instruction.Operation,
			Left:  int64(*left),
			Right: int64(*right),
		}

		fmt.Print("start ")
		response, err := s.client.Calculate(context.Background(), &req)
		fmt.Println("end")
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

func (s *Service) GetDependentValues(key string) *[]Instruction {
	return s.instructionStorage.Get(key)
}

func (s *Service) DeleteKey(key string) {
	s.instructionStorage.Delete(key)
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
