package services

type ArithmeticService struct {
}

func NewArithmeticService() *ArithmeticService {
	return &ArithmeticService{}
}

func (as *ArithmeticService) Sum(left, right int64) int64 {
	return left + right
}

func (as *ArithmeticService) Multi(left, right int64) int64 {
	return left * right
}

func (as *ArithmeticService) Sub(left, right int64) int64 {
	return left - right
}
