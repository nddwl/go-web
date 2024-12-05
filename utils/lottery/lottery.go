package lottery

import (
	"math/rand"
	"sort"
	"strconv"
	"time"
)

var (
	rng = &rand.Rand{}
)

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func Draw() (int, string) {
	prize := []int{10, 25, 50, 100}
	cumulativeWeights := []float64{0.5, 0.75, 0.9, 1}
	randomWeight := rng.Float64() * cumulativeWeights[len(cumulativeWeights)-1]
	index := sort.Search(len(cumulativeWeights), func(i int) bool {
		return cumulativeWeights[i] >= randomWeight
	})
	return prize[index], "签到成功,恭喜抽到" + strconv.Itoa(prize[index]) + "硬币!"
}
