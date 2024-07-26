package contoroller

import (
	"context"
	"tenkiej/config"
	"tenkiej/contoroller/amedas"
	"tenkiej/logger"
	"time"
)

var shutdown chan struct{}
var done chan struct{}
var wait chan struct{}

var amedasAPI *amedas.API

const (
	WAITTIME_DEFAULT = 1 * time.Minute
)

func Init() {
	shutdown = make(chan struct{}, 1)
	done = make(chan struct{}, 1)
	wait = make(chan struct{}, 1)
	amedasAPI = amedas.NewAmedasAPI()
}

func Run(ctx context.Context) {
	var oneshut chan struct{} = make(chan struct{}, 1)
	oneshut <- struct{}{}
	waitTime := WAITTIME_DEFAULT
	var TabledataGetTime time.Time
	logger.Info("Start Run")
loop:
	for {
		select {
		case <-shutdown:
			done <- struct{}{}
			break loop
		case <-oneshut:
			ctx1, span := config.TracerS(ctx, "controller.Run", "oneshut")
			stime := time.Now()
			amedasAPI.GetAmedasMapData(ctx1)
			amedasAPI.GetAmedasTableData(ctx1)
			_, span1 := config.TracerS(ctx1, "controller.Run", "GetPrometesusData")
			amedasAPI.GetPrometesusData()
			span1.End()
			_, span1 = config.TracerS(ctx1, "controller.Run", "GetPrometesusDatas")
			amedasAPI.GetPrometesusDatas()
			amedasAPI.CleanData()
			span1.End()
			TabledataGetTime = stime
			logger.Info("1 time after loop",
				"Time", time.Since(stime).String(),
				"DataCount", len(amedasAPI.PData),
				"DataMapCount", amedasAPI.CountSumMap(),
			)
			waitTime = WAITTIME_DEFAULT - time.Now().Sub(stime)
			span.End()
		case <-wait:
			done <- struct{}{}
		case <-time.After(waitTime):
			ctx1, span := config.TracerS(ctx, "controller.Run", "1 time after")
			stime := time.Now()
			amedasAPI.GetAmedasMapData(ctx1)
			//60分以上経過していたらTableDataも取得
			if time.Since(TabledataGetTime) > 60*time.Minute {
				amedasAPI.GetAmedasTableData(ctx1)
				TabledataGetTime = stime
			}
			_, span1 := config.TracerS(ctx1, "controller.Run", "GetPrometesusData")
			amedasAPI.GetPrometesusData()
			span1.End()
			_, span1 = config.TracerS(ctx1, "controller.Run", "GetPrometesusDatas")
			amedasAPI.GetPrometesusDatas()
			amedasAPI.CleanData()
			span1.End()
			logger.Info("1 time after loop",
				"Time", time.Since(stime).String(),
				"DataCount", len(amedasAPI.PData),
				"DataMapCount", amedasAPI.CountSumMap(),
			)
			waitTime = WAITTIME_DEFAULT - time.Now().Sub(stime)
			span.End()
			// default:
		}
	}
	logger.Info("End Run")
}

func Wait(ctx context.Context) error {
	wait <- struct{}{}
	ctx, stop := context.WithTimeout(ctx, 10*time.Minute)
	defer stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		break
	}
	return nil
}

func Stop(ctx context.Context) {
	shutdown <- struct{}{}
	ctx, stop := context.WithTimeout(ctx, 1*time.Minute)
	defer stop()
	select {
	case <-ctx.Done():
		logger.Error("Stop", "Error", ctx.Err())
		break
	case <-done:
		break
	}
	close(shutdown)
	close(done)
}
