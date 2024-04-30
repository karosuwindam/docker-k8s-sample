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
	if err := senser.Run(ctx); err != nil {
		return errors.Wrap(err, "senser.Run()")
	}
	return nil
}

func Stop(ctx context.Context) error {
	if err := senser.Stop(ctx); err != nil {
		return err
	}
	return nil
}

func Wait() {
	senser.Wait()
}
