package movo

const (
    minute int64 = 60
    hour int64 = 60*minute
    day int64 = 24*hour
    )

func MinuteToSeconds(m int) int64 {
  return minute*int64(m)
}

func HourToSeconds(h int) int64 {
  return hour*int64(h)
}

func DayToSeconds(d int) int64 {
  return day*int64(d)
}
