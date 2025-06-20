package middleware

import (
	"fmt"
	"sync"
	"sync/atomic"
	"wellnesspath/config"

	"github.com/gin-gonic/gin"
)

var (
	globalConcurrentRequests int32
	queueThreshold           int32 = 30
	globalQueue                    = make(chan func(), 1000)
	once                     sync.Once
)

func QueueMiddleware() gin.HandlerFunc {
	once.Do(func() {
		go func() {
			for task := range globalQueue {
				if config.ENV.Queue == "concurrent" {
					go task() // run in parallel
				} else {
					task() // run sequentially
				}
			}
		}()
	})

	return func(c *gin.Context) {
		atomic.AddInt32(&globalConcurrentRequests, 1)
		defer atomic.AddInt32(&globalConcurrentRequests, -1)

		curr := atomic.LoadInt32(&globalConcurrentRequests)

		if curr >= queueThreshold {
			fmt.Printf("⚠️  QUEUE TRIGGERED (%s) — Current: %d\n", config.ENV.Queue, curr)

			done := make(chan struct{})

			globalQueue <- func() {
				fmt.Println("🔁 Running task from queue...")
				c.Next()
				close(done)
			}

			<-done // wait until the task completes to send response
		} else {
			fmt.Printf("✅ Direct handling — Current: %d\n", curr)
			c.Next()
		}
	}
}
