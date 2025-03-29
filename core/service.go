package core

import "errors"

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s *Service) Calculate(op string, left, right int) (*int, error) {
	var res int
	switch op {
	case "+":
		res = s.Sum(left, right)
	case "*":
		res = s.Multi(left, right)
	case "-":
		res = s.Sub(left, right)
	default:
		return nil, errors.New("unknown operation")
	}
	return &res, nil
}

func (s *Service) Sum(left, right int) int {
	return left + right
}

func (s *Service) Multi(left, right int) int {
	return left * right
}

func (s *Service) Sub(left, right int) int {
	return left - right
}
