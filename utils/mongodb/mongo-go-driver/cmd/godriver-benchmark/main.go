package main

import (
	"os"

	"GDTS/utils/mongodb/mongo-go-driver/benchmark"
)

func main() {
	os.Exit(benchmark.DriverBenchmarkMain())
}
