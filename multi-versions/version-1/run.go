package main

import (
	"context"
	"github.com/aionbit/ixo-plugins/multi-versions/lib"
	"github.com/aionbit/ixo-plugins/multi-versions/version-1/common"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

type request any

// Input 输入
func Input() interface{} {
	var r request
	return &r
}

// Run Note
func Run(ctx context.Context, input interface{}) (interface{}, error) {
	data, err := oj.ParseString(lib.Data())
	if err != nil {
		return nil, err
	}
	setExpr, err := jp.ParseString(common.Path())
	if err != nil {
		return nil, err
	}
	if err = setExpr.Set(data, "v1.26.7"); err != nil {
		return nil, err
	}
	return oj.JSON(data), nil
}
