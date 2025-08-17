package main

import (
	"log"

	"github.com/go-redis/redis/v8"

	"api-gateway-go/internal/jobs"
	"api-gateway-go/internal/kafka"
	"api-gateway-go/internal/loadbalancer"
	"api-gateway-go/internal/metrics"
	"api-gateway-go/internal/middleware"
	"api-gateway-go/internal/server"
)

func main() {
	// Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	mw := &middleware.RedisMiddleware{Rdb: rdb}

	// Load balancer
	lb := loadbalancer.NewRoundRobin([]string{
		"http://service1-instance1:8081",
		"http://service1-instance2:8082",
	})

	// Job queue
	jobQueue := make(chan jobs.Job, 100)
	jobs.StartWorkers(jobQueue, 10)

	// Kafka
	producer := kafka.NewProducer([]string{"localhost:9092"}, "my-topic")
	defer producer.Close()
	producer.Publish("Service started")

	// Start Prometheus metrics
	metrics.StartMetricsServer(":9000")

	// Start Fiber server
	app := server.NewServer(mw, lb)
	log.Fatal(app.Listen(":8080"))
}
