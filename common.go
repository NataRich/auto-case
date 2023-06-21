package main

import (
  "os"
  "fmt"
  "log"
  "time"
  "math/rand"
)

func NowYearStr() string {
  return fmt.Sprintf("%d", time.Now().Year())
}

func DeltaDayStr(delta int, base time.Time) string {
  base = base.AddDate(0, 0, delta)
  return string(base.Format("2006-01-02 15:04:05"))
}

func InsertRandomDates(conf *CaseConfig) {
  delta := rand.Intn(10) + 1
  endDelta := rand.Intn(10) + 1
  
  now := time.Now()
  conf.StartTime = DeltaDayStr(endDelta + delta, now)
  conf.EndTime = DeltaDayStr(endDelta, now)

  DebugPrint(fmt.Sprintf("起始时间为：%s\n", conf.StartTime))
  DebugPrint(fmt.Sprintf("结束时间为：%s\n", conf.EndTime))
}

func LogError(err error, path string) {
  f, er := os.OpenFile(path, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
  if er != nil {
    log.Println("无法打开日志文件\n")
  }

  defer f.Close()

  if _, er := f.WriteString(fmt.Sprintf("%v\n", err)); er != nil {
    log.Println("无法写入日志文件\n")
  }
}
