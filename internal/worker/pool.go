// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package worker // import "miniflux.app/v2/internal/worker"

import (
	"log/slog"

	"miniflux.app/v2/internal/model"
	"miniflux.app/v2/internal/storage"
)

// Pool handles a pool of workers.
type Pool struct {
	queue chan model.Job
}

// Push send a list of jobs to the queue.
func (p *Pool) Push(jobs model.JobList) {
	nbJobs := len(jobs)
	if nbJobs > 0 {
		slog.Info("Created a batch of feeds",
			slog.Int("nb_jobs", nbJobs),
		)
	}
	for _, job := range jobs {
		p.queue <- job
	}
}

// NewPool creates a pool of background workers.
func NewPool(store *storage.Storage, nbWorkers int) *Pool {
	workerPool := &Pool{
		queue: make(chan model.Job),
	}

	for i := 0; i < nbWorkers; i++ {
		worker := &Worker{id: i, store: store}
		go worker.Run(workerPool.queue)
	}

	return workerPool
}
