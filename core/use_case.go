package core

import (
	"errors"
	"strconv"
)

type UseCase struct {
	service Service
}

func NewUseCase(s Service) UseCase {
	return UseCase{service: s}
}

func (uc *UseCase) Execute(dtos []DTO) ([]Item, error) {

	buff := make(map[string]int)
	var items []Item
	var printFuncs []DTO

	dep := make(map[string][]string)
	var waiters []DTO

	for _, dto := range dtos {
		if dto.Type == "calc" {
			if _, ok := buff[dto.Result]; ok {
				return nil, errors.New("var already used")
			}

			left, err := strconv.Atoi(*dto.Left)
			if err != nil {
				if _, ok := buff[*dto.Left]; !ok {
					//return nil, errors.New("unknown left op")

					dep[*dto.Left] = append(dep[*dto.Left], dto.Result)
					waiters = append(waiters, dto)
					continue
				}

				left = buff[*dto.Left]
			}

			right, err := strconv.Atoi(*dto.Right)
			if err != nil {
				if _, ok := buff[*dto.Right]; !ok {
					//return nil, errors.New("unknown right op")

					dep[*dto.Right] = append(dep[*dto.Right], dto.Result)
					waiters = append(waiters, dto)
					continue
				}
				right = buff[*dto.Right]
			}

			res, err := uc.service.Calculate(*dto.Operation, left, right)

			if err != nil {
				return nil, err
			}

			buff[dto.Result] = *res

		} else if dto.Type == "print" {
			printFuncs = append(printFuncs, dto)
		} else {
			return nil, errors.New("unknown type")
		}
	}

	for _, waiter := range waiters {
		left, err := strconv.Atoi(*waiter.Left)
		if err != nil {
			if _, ok := buff[*waiter.Left]; !ok {
				return nil, errors.New("unknown left op")
			}

			left = buff[*waiter.Left]
		}

		right, err := strconv.Atoi(*waiter.Right)
		if err != nil {
			if _, ok := buff[*waiter.Right]; !ok {
				return nil, errors.New("unknown right op")
			}
			right = buff[*waiter.Right]
		}

		res, err := uc.service.Calculate(*waiter.Operation, left, right)

		if err != nil {
			return nil, err
		}

		buff[waiter.Result] = *res
	}

	for _, printFunc := range printFuncs {
		items = append(items, Item{Var: printFunc.Result, Value: buff[printFunc.Result]})
	}
	return items, nil
}
