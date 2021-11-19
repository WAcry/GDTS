package main

import (
	"os"

	"utils/mongodb/mongo-go-driver/benchmark"
)

func main() {
	os.Exit(benchmark.DriverBenchmarkMain())
}
