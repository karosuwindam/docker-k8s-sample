package common

import "context"

func Run(ctx context.Context, ch chan<- error) {}

func Setup() error {
	return nil
}
