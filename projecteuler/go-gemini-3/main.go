package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type Result struct {
	Problem int     `json:"solution"`
	Answer  string  `json:"result"`
	TimeMs  float64 `json:"time_ms"`
}

var solutions = map[int]func() string{
	1: solution_0001,
	2: solution_0002,
	3: solution_0003,
	4: solution_0004,
	5: solution_0005,
	6: solution_0006,
	7: solution_0007,
	8: solution_0008,
	9: solution_0009,
	10: solution_0010,
	350: solution_0350,
	387: solution_0387,
	407: solution_0407,
}

func main() {
	var keys []int
	for k := range solutions {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		start := time.Now()
		ans := solutions[k]()
		elapsed := float64(time.Since(start).Microseconds()) / 1000.0

		res := Result{
			Problem: k,
			Answer:  ans,
			TimeMs:  elapsed,
		}
		
		data, _ := json.Marshal(res)
		fmt.Println(string(data))
	}
}
