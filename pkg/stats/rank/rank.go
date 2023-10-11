// Copyright 2023 PraserX
package rank

// RankCnt defines total number of ranks.
const RankCnt = 8

// Ranks definition.
var Ranks = [RankCnt]string{
	"Unranked",
	"Coffee trainee",
	"Hesitant customer",
	"Regular customer",
	"Coffee business supporter",
	"Coffee gourmet",
	"Coffee guru",
	"God-Emperor",
}

// RankLimits specifies required number of consumed coffees to acquired given
// rank.
var RankLimits = [RankCnt]int{
	0,
	3,
	8,
	12,
	18,
	25,
	35,
	45,
}

// ComputeRank computes and return rank based on monthly consumed coffee cups.
func ComputeRank(qpm int) (userRank string) {
	userRank = Ranks[0]
	for i := 0; i < RankCnt; i++ {
		if RankLimits[i] <= qpm {
			userRank = Ranks[i]
		}
	}
	return userRank
}
