package goutils

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileSize(t *testing.T) {
	ta := assert.New(t)
	ta.Equal("0 B", FileSize(0))
	ta.Equal("1023 B", FileSize(1023))
	ta.Equal("1.0 KB", FileSize(1024))
	ta.Equal("1.0 MB", FileSize(1024*1024))
}

func TestMakeRange(t *testing.T) {
	ta := assert.New(t)

	ta.Nil(MakeRange(2, 1))
	ta.Equal([]int{1, 2, 3, 4, 5}, MakeRange(1, 5))
	ta.Equal([]int{2}, MakeRange(2, 2))
	ta.Equal([]int{-3, -2, -1, 0, 1, 2}, MakeRange(-3, 2))
}

func TestRound(t *testing.T) {
	var vf = []struct {
		v                              float64
		pm1                            float64
		p0, p1, p2, p3, p4, p5, p6, p7 float64
	}{
		{4.9790119248836735e+00, 0, 5, 5.0, 4.98, 4.979, 4.9790, 4.97901, 4.979012, 4.9790119},
		{7.7388724745781045e+00, 10, 8, 7.7, 7.74, 7.739, 7.7389, 7.73887, 7.738872, 7.7388725},
		{-2.7688005719200159e-01, 0, 0, -0.3, -0.28, -0.277, -0.2769, -0.27688, -0.276880, -0.2768801},
		{-5.0106036182710749e+00, -10, -5, -5.0, -5.01, -5.011, -5.0106, -5.01060, -5.010604, -5.0106036},
		{9.6362937071984173e+00, 10, 10, 9.6, 9.64, 9.636, 9.6363, 9.63629, 9.636294, 9.6362937},
		{2.9263772392439646e+00, 0, 3, 2.9, 2.93, 2.926, 2.9264, 2.92638, 2.926377, 2.9263772},
		{5.2290834314593066e+00, 10, 5, 5.2, 5.23, 5.229, 5.2291, 5.22908, 5.229083, 5.2290834},
		{2.7279399104360102e+00, 0, 3, 2.7, 2.73, 2.728, 2.7279, 2.72794, 2.727940, 2.7279399},
		{1.8253080916808550e+00, 0, 2, 1.8, 1.83, 1.825, 1.8253, 1.82531, 1.825308, 1.8253081},
		{-8.6859247685756013e+00, -10, -9, -8.7, -8.69, -8.686, -8.6859, -8.68592, -8.685925, -8.6859248},

		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1.390671161567e-309, 0, 0, 0, 0, 0, 0, 0, 0, 0},               // denormal
		{0.49999999999999994, 0, 1, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5}, // 0.5-epsilon
		{0.5, 0, 1, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
		{0.5000000000000001, 0, 1, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5}, // 0.5+epsilon
		{-1.5, 0, -1, -1.5, -1.5, -1.5, -1.5, -1.5, -1.5, -1.5},
		{-2.5, 0, -2, -2.5, -2.5, -2.5, -2.5, -2.5, -2.5, -2.5},
		{math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN()},
		{math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1)},
	}

	var vfSC = []struct {
		v                                 float64
		p1, p0                            float64
		pm1, pm2, pm3, pm4, pm5, pm6, pm7 float64
	}{
		{2251799813685249.5, 2251799813685249.5, 2251799813685250, 2251799813685250, 2251799813685200, 2251799813685000, 2251799813690000, 2251799813700000, 2251799814000000, 2251799810000000},
		{2251799813685250.5, 2251799813685250.5, 2251799813685251, 2251799813685250, 2251799813685300, 2251799813685000, 2251799813690000, 2251799813700000, 2251799814000000, 2251799810000000},
		{4503599627370495.5, 4503599627370495, 4503599627370496, 4503599627370500, 4503599627370500, 4503599627370000, 4503599627370000, 4503599627400000, 4503599627000000, 4503599630000000},
		{4503599627370497, 4503599627370497, 4503599627370498, 4503599627370500, 4503599627370500, 4503599627370000, 4503599627370000, 4503599627400000, 4503599627000000, 4503599630000000},
	}

	for _, cases := range vf {
		expected := []float64{cases.pm1, cases.p0, cases.p1, cases.p2, cases.p3, cases.p4, cases.p5, cases.p6, cases.p7}
		for precision := -1; precision < 8; precision++ {
			if actual := Round(cases.v, precision); actual != expected[precision+1] && !(math.IsNaN(actual) && math.IsNaN(expected[precision+1])) {
				t.Errorf("Round(%f, %d) %f != %f", cases.v, precision, actual, expected[precision+1])
			}
		}
	}

	for _, cases := range vfSC {
		expected := []float64{cases.p1, cases.p0, cases.pm1, cases.pm2, cases.pm3, cases.pm4, cases.pm5, cases.pm6, cases.pm7}
		for i, precision := range []int{1, 0, -1, -2, -3, -4, -5, -6, -7} {
			if actual := Round(cases.v, precision); actual != expected[i] {
				t.Errorf("Round(%f, %d) %f != %f", cases.v, precision, actual, expected[i])
			}
		}
	}

	ta := assert.New(t)
	ta.Equal(0.364, Round(0.363636, 3))
	ta.Equal(0.36, Round(0.363636, 2))
	ta.Equal(0.4, Round(0.363636, 1))
	ta.Equal(30.0, Round(32, -1))
}
