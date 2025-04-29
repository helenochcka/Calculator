package services

type ArithmService struct {
}

func NewArithmService() ArithmService {
	return ArithmService{}
}

func (as *ArithmService) Sum(left, right int64) int64 {
	return left + right
}

func (as *ArithmService) Multi(left, right int64) int64 {
	return left * right
}

func (as *ArithmService) Sub(left, right int64) int64 {
	return left - right
}
