package friendlycaptcha

import (
	"runtime/debug"
)

func Version() string {
	if bi, ok := debug.ReadBuildInfo(); ok {
		return bi.Main.Version
	}
	return "unknown"
}
