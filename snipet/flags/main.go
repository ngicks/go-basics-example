package main

import (
	"flag"
	"fmt"
	"log/slog"
	"time"
)

var (
	flag1        = flag.String("f1", "", "flag 1")
	flag2        = flag.Bool("f2", false, "flag 2")
	flag3        = flag.Int("f3", 0, "flag 3")
	flagMultiple string
)

func init() {
	flag.StringVar(&flagMultiple, "fm1", "1", "flag multiple 1")
	flag.StringVar(&flagMultiple, "fm2", "2", "flag multiple 2")
	flag.StringVar(&flagMultiple, "fm3", "3", "flag multiple 3")
}

var (
	t1       time.Time
	logLevel slog.Level = 99999
	t2       time.Time
)

func init() {
	flag.Func("t1", "time to start", func(s string) error {
		var err error
		t1, err = time.Parse(time.RFC3339, s)
		return err
	})
	flag.BoolFunc("log", "bool func", func(s string) error {
		switch s {
		case "true":
			logLevel = slog.LevelInfo
		case "":
		default:
			return logLevel.UnmarshalText([]byte(s))
		}
		return nil
	})
}

func init() {
	flag.TextVar(&t2, "t2", time.Now(), "time to end")
}

func main() {
	flag.Parse()

	for _, s := range [][2]any{
		{"flag1", *flag1},
		{"flag2", *flag2},
		{"flag3", *flag3},
		{"flagMultiple", flagMultiple},
		{"t1", t1},
		{"t2", t2},
		{"log", logLevel},
	} {
		fmt.Printf("%s = %v\n", s[0], s[1])
	}

	for i := range 3 {
		positionalArg := flag.Arg(i)
		fmt.Printf("position %d = %s\n", i, positionalArg)
	}
	fmt.Printf("args = %#v\n", flag.Args())
}
