package main

import (
	"fmt"
	"time"
)

func main() {
	day := time.Now()
	const layout = "2006-01-02"
	fmt.Println(day.Format(layout))
}
