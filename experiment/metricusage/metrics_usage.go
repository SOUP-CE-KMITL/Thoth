package main

import (
     "fmt"
     "time"
     _ "github.com/shirou/gopsutil/mem"
     "github.com/shirou/gopsutil/cpu"
)

func main() {
     //v, _ := mem.VirtualMemory()
     cpu_percent, _ := cpu.CPUPercent(time.Duration(1) * time.Second, false)

     // almost every return value is a struct
     // fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

     fmt.Println(cpu_percent)
     // convert to JSON. String() is also implemented
     //fmt.Println(v)
}