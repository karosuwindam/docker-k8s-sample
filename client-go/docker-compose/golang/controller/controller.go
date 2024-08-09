package controller

import (
	"context"
	"errors"
	"ingresslist/controller/datastore"
	"ingresslist/controller/kubeget"
	"log/slog"
	"time"
)

var shutdown chan struct{}
var done chan struct{}

func Init() error {
	if err := kubeget.Init(); err != nil {
		return err
	}
	if err := datastore.Init(); err != nil {
		return err
	}
	shutdown = make(chan struct{}, 1)
	done = make(chan struct{}, 1)
	return nil
}

func Run(ctx context.Context) error {
	oneshut := make(chan struct{}, 1)
	oneshut <- struct{}{}
	slog.Info("Controller Run Start")
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-shutdown:
			done <- struct{}{}
			break loop
		case <-oneshut:
			k8sGet(ctx)
			done <- struct{}{}
			//
		case <-time.After(10 * time.Second):
			k8sGet(ctx)
			//
		}
	}
	slog.Info("Controller Run Stop")
	return nil
}

func Stop(ctx context.Context) error {
	done = make(chan struct{}, 1)
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	shutdown <- struct{}{}
	select {
	case <-ctx.Done():
		return errors.New("Stop Time out")
	case <-done:
		break
	}
	return nil
}

func Wait() {
	select {
	case <-done:
		break
	case <-time.After(time.Second):
		slog.Error("Wait time out 1 sec")
		break
	}
}

func k8sGet(ctx context.Context) {
	var startCh chan bool = make(chan bool, 5)
	var endCh chan bool = make(chan bool, 5)
	startCh <- true
	go func(ctx context.Context) {
		defer func() {
			endCh <- <-startCh
		}()
		if node, err := kubeget.GetNodeData(ctx); err != nil {
			slog.Error("node get", "error", err)
		} else {
			datastore.Write(node)
		}
	}(ctx)
	startCh <- true
	go func(ctx context.Context) {
		defer func() {
			endCh <- <-startCh
		}()
		if pod, err := kubeget.GetPodData(ctx); err != nil {
			slog.Error("pod get", "error", err)
		} else {
			datastore.Write(pod)
		}
	}(ctx)
	startCh <- true
	go func(ctx context.Context) {
		defer func() {
			endCh <- <-startCh
		}()
		if ingress, err := kubeget.GetIngressData(ctx); err != nil {
			slog.Error("ingress get", "error", err)
		} else {
			datastore.Write(ingress)
		}
	}(ctx)
	startCh <- true
	go func(ctx context.Context) {
		defer func() {
			endCh <- <-startCh
		}()
		if service, err := kubeget.GetServiceData(ctx); err != nil {
			slog.Error("service get", "error", err)
		} else {
			datastore.Write(service)
		}
	}(ctx)
loop:
	for {
		select {
		case <-ctx.Done():
			return
		case <-endCh:
			if len(startCh) == 0 {
				break loop
			}
		}
	}
	close(startCh)
	close(endCh)
	return
}
