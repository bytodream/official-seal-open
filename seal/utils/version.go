package utils

import "time"

// Version is the library version.
// 声明为变量以便构建时通过 -ldflags "-X ...utils.Version=vX.Y.Z" 注入（见 script/build.sh）。
var Version = "1.0.0"

// ReleasedAt is the release time of the current library version.
var ReleasedAt = time.Date(2026, 11, 16, 15, 30, 0, 0, time.UTC)
