package main

import (
	"context"
	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/redisq"
	"log"
)

var MainQueue taskq.Queue
var CompressTask *taskq.Task
var libvipsCompressor = LibVipsCompressor{}
var lilliputCompressor = LilliputCompressor{}

func PrepareQueueHandler() error {
	QueueFactory := redisq.NewFactory()

	MainQueue = QueueFactory.RegisterQueue(&taskq.QueueOptions{
		Name:          "api-worker",
		Redis:         Cache.rdb,
		MinNumWorker:  1,
		MaxNumWorker:  5,
		MaxNumFetcher: 5,
	})

	CompressTask = taskq.RegisterTask(&taskq.TaskOptions{
		Name: "compress",
		Handler: func(id, compressor string, quality, progressive int) error {

			imageOpts := ImageOpts{
				id:          id,
				compressor:  compressor,
				quality:     quality,
				progressive: progressive,
			}

			var c Compressor
			if imageOpts.compressor == "libvips" {
				c = libvipsCompressor
			} else {
				c = lilliputCompressor
			}

			Image{
				id:          imageOpts.id,
				compressor:  c,
				Quality:     imageOpts.quality,
				Progressive: imageOpts.progressive,
			}.Compress()
			return nil
		},
	})

	consumer := taskq.NewConsumer(MainQueue)

	return consumer.Start(context.Background())
}

func QueueJob(id, compressor string, quality, progressive int) {
	ctx := context.Background()
	msg := CompressTask.WithArgs(ctx, id, compressor, quality, progressive)
	msg.Delay = 0

	err := MainQueue.Add(msg)
	if err != nil {
		log.Println(err)
	}
}
