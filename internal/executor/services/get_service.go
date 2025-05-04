package services

import (
	"Calculator/internal/executor"
	"Calculator/internal/executor/values"
	"context"
)

type GetterService struct {
}

func NewGetterService() *GetterService {
	return &GetterService{}
}

func (gs *GetterService) GetReqIdFromCtx(ctx context.Context) (*string, error) {
	reqId, exists := ctx.Value(values.RequestIdKey).(string)
	if !exists {
		return nil, executor.ErrReqIdMissing
	}
	return &reqId, nil
}

func (gs *GetterService) GetVarValue(variable interface{}, resultMap map[string]int) (*int, bool) {
	switch variable.(type) {
	case string:
		res, ok := resultMap[variable.(string)]
		if !ok {
			return nil, false
		}
		return &res, true
	case int:
		value := variable.(int)
		return &value, true
	}
	return nil, false
}
