// A job runner for executing scheduled or ad-hoc tasks asynchronously from HTTP requests.
//
// It adds a couple of features on top of the Robfig cron package:
// 1. Protection against job panics.  (They print to ERROR instead of take down the process)
// 2. (Optional) Limit on the number of jobs that may run simulatenously, to
//    limit resource consumption.
// 3. (Optional) Protection against multiple instances of a single job running
//    concurrently.  If one execution runs into the next, the next will be queued.
// 4. Cron expressions may be defined in app.conf and are reusable across jobs.
// 5. Job status reporting. [WIP]
package jobrunner

import (
	"time"

	"github.com/robfig/cron/v3"
)

// Callers can use jobs.Func to wrap a raw func.
// (Copying the type to this package makes it more visible)
//
// For example:
//    jobrunner.Schedule("cron.frequent", jobs.Func(myFunc))
type Func func()

func (r Func) Run() { r() }

// spec 时间设置
// job 任务配置
// des 任务描述
func Schedule(spec string, job cron.Job, s ...string) error {
	sched, err := cron.ParseStandard(spec)
	if err != nil {
		return err
	}
	MainCron.Schedule(sched, New(job, s...))
	return nil
}

// Run the given job at a fixed interval.
// The interval provided is the time between the job ending and the job being run again.
// The time that the job takes to run is not included in the interval.
func Every(duration time.Duration, job cron.Job, s ...string) {

	MainCron.Schedule(cron.Every(duration), New(job, s...))
}

// Run the given job right now.
func Now(job cron.Job, s ...string) {
	go New(job, s...).Run()
}

// Run the given job once, after the given delay.
func In(duration time.Duration, job cron.Job, s ...string) {
	go func() {
		time.Sleep(duration)
		New(job, s...).Run()
	}()
}
