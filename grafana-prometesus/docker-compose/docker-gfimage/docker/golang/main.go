package main

import (
	"math"

	"gonum.org/v1/gonum/graph/simple"
)

func main() {
	self := 0             // the cost of self connection
	absent := math.Inf(1) // the wieght returned for absent edges

	graph := simple.NewWeightedUndirectedGraph(self, absent)
}
