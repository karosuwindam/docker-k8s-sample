package controller

import (
	"context"
	"eijent/controller/senser"

	"github.com/pkg/errors"
)

func Init() error {
	if err := senser.Init(); err != nil {
		return errors.Wrap(err, "senser.Init()")
	}
	return nil
}

func Run(ctx context.Context) error {
	return nil
}

func Stop() error {
	return nil
}
