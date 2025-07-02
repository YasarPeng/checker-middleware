package util

import (
	"fmt"

	"github.com/spf13/pflag"
)

func PrintFlags(flags *pflag.FlagSet) {
	//
	flags.VisitAll(func(f *pflag.Flag) {
		PrintFlag(f)
	})
}

func PrintFlag(f *pflag.Flag) {
	//
	shorthandPart := ""
	if f.Shorthand != "" {
		shorthandPart = fmt.Sprintf("-%s, ", f.Shorthand)
	}

	// 构建 default 部分
	defaultPart := ""
	if f.DefValue != "" && f.DefValue != "[]" && f.DefValue != "false" {
		defaultPart = fmt.Sprintf(" (default: %v)", f.DefValue)
	}

	fmt.Printf("%5s--%-10s %s%s\n", shorthandPart, f.Name, f.Usage, defaultPart)
}
