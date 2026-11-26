package base

import "testing"

// TestNewRandFixedSeedReproducible 固定种子时，两次取源应产生相同序列。
func TestNewRandFixedSeedReproducible(t *testing.T) {
	t.Setenv(RandomSeedEnv, "42")
	a := NewRand()
	b := NewRand()
	for i := 0; i < 20; i++ {
		if got, want := a.Int63(), b.Int63(); got != want {
			t.Fatalf("第 %d 个随机数不一致: %d != %d", i, got, want)
		}
	}
}

// TestNewRandInvalidSeedFallsBack 非法或 0 种子按随机处理，两次取源不应相同。
func TestNewRandInvalidSeedFallsBack(t *testing.T) {
	t.Setenv(RandomSeedEnv, "0")
	a := NewRand().Int63()
	b := NewRand().Int63()
	if a == b {
		t.Fatalf("0 种子应回退为随机，但两次结果相同: %d", a)
	}

	t.Setenv(RandomSeedEnv, "not-a-number")
	a = NewRand().Int63()
	b = NewRand().Int63()
	if a == b {
		t.Fatalf("非法种子应回退为随机，但两次结果相同: %d", a)
	}
}
