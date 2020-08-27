package main

import (
	"context"
	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/redisq"
	"log"
)

var MainQueue taskq.Queue
var CompressTask *taskq.Task

func PrepareQueueHandler() error {
	QueueFactory := redisq.NewFactory()

	MainQueue = QueueFactory.RegisterQueue(&taskq.QueueOptions{
		Name:        "api-worker",
		Redis:       cache.rdb,
		MinNumWorker: 1,
		MaxNumWorker: 5,
		MaxNumFetcher: 5,
	})

	CompressTask = taskq.RegisterTask(&taskq.TaskOptions{
		Name: "compress",
		Handler: func(id string) error {
			Image{id: id}.Compress()
			return nil
		},
	})

	consumer := taskq.NewConsumer(MainQueue)

	return consumer.Start(context.Background())
}

func QueueJob(id string) {
	ctx := context.Background()
	msg := CompressTask.WithArgs(ctx, id)
	msg.Delay = 0
	err := MainQueue.Add(msg)
	if err != nil {
		log.Fatal(err)
	}
}
