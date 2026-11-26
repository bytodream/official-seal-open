package base

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// RandomSeedEnv 固定随机种子的环境变量名。
// 设为非 0 整数时，防伪断裂线与模糊加粗的随机绘制结果可复现（便于图像回归比对）；
// 未设置或为 0 时按时间取种子，保持每次生成略有差异的现有行为。
const RandomSeedEnv = "OFFICIAL_SEAL_RANDOM_SEED"

// fixedRandomSeed 读取并解析环境变量，返回种子与是否启用固定种子。
// 每次调用都读取环境变量，便于调用方在进程内（如命令行参数解析后）即时生效。
func fixedRandomSeed() (int64, bool) {
	raw := strings.TrimSpace(os.Getenv(RandomSeedEnv))
	if raw == "" {
		return 0, false
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || v == 0 {
		return 0, false
	}
	return v, true
}

// NewRand 返回绘制用随机源。
// 设置了 OFFICIAL_SEAL_RANDOM_SEED 时每次调用都从同一起点展开，
// 因此相同入参的两次渲染结果一致；否则以当前时间作为种子。
func NewRand() *rand.Rand {
	if seed, ok := fixedRandomSeed(); ok {
		return rand.New(rand.NewSource(seed))
	}
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
