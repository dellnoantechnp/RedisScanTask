package utils

import (
	"github.com/fatih/color"
)

// @Title        branchPrefix.go
// @Description
// @Create       2026-04-16 17:38
// @Update       2026-04-16 17:38

var (
	Colorize     bool
	branchPrefix = "\\__"
)

func ColorizePrefix() string {
	if !Colorize {
		return color.HiBlackString(branchPrefix)
	} else {
		return branchPrefix
	}
}
