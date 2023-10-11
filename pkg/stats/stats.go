package stats

import (
	"errors"
	"fmt"

	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/models"
	"github.com/praserx/gobarista/pkg/stats/rank"
	"gorm.io/gorm"
)

type Stats struct {
	User   UserStats
	Period PeriodStats
}

type UserStats struct {
	Rank                   string
	CurrentAverageCoffees  int
	PrevAverageCoffees     int
	PrevTotalCoffees       int
	PrevTotalMonths        int
	CoffeeConsumptionTrend string
}

type PeriodStats struct {
	TotalCustomers int
	TotalMonths    int
	TotalQuantity  int
	TotalAverage   int
}

func GetStats(uid uint, currentPeriod models.Period, currentBill models.Bill, totalCustomers int) (stats Stats, err error) {
	userStats, err := GetUserStats(uid, currentPeriod, currentBill)

	return Stats{
		User:   userStats,
		Period: GetPeriodStats(currentPeriod, totalCustomers),
	}, err
}

func GetPeriodStats(currentPeriod models.Period, totalCustomers int) (stats PeriodStats) {
	return PeriodStats{
		TotalCustomers: totalCustomers,
		TotalMonths:    GetTotalMonths(currentPeriod),
		TotalQuantity:  currentPeriod.TotalQuantity,
		TotalAverage:   currentPeriod.TotalQuantity / GetTotalMonths(currentPeriod),
	}
}

func GetUserStats(uid uint, currentPeriod models.Period, currentBill models.Bill) (stats UserStats, err error) {
	prevBills, err := database.SelectAllBillsForUserExceptSpecifiedPeriod(uid, currentPeriod.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return UserStats{}, fmt.Errorf("billing period not found")
	} else if err != nil {
		return UserStats{}, fmt.Errorf("cannot get billing period: %v", err)
	}

	prevPeriods, err := database.SelectAllPeriodsExceptSpecifiedPeriod(currentPeriod.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return UserStats{}, fmt.Errorf("billing period not found")
	} else if err != nil {
		return UserStats{}, fmt.Errorf("cannot get billing period: %v", err)
	}

	for _, prevBill := range prevBills {
		stats.PrevTotalCoffees += prevBill.Quantity
	}

	for _, prevPeriod := range prevPeriods {
		stats.PrevTotalMonths += GetTotalMonths(prevPeriod)
	}

	stats.PrevAverageCoffees = stats.PrevTotalCoffees / stats.PrevTotalMonths
	stats.CurrentAverageCoffees = (stats.PrevTotalCoffees + currentBill.Quantity) / (stats.PrevTotalMonths + GetTotalMonths(currentPeriod))

	if stats.PrevAverageCoffees > stats.CurrentAverageCoffees {
		stats.CoffeeConsumptionTrend = "<span style=\"color: #cb4335;\">&#129126;</span>"
	} else if stats.PrevAverageCoffees < stats.CurrentAverageCoffees {
		stats.CoffeeConsumptionTrend = "<span style=\"color: #2ecc71;\">&#129125;</span>"
	} else {
		stats.CoffeeConsumptionTrend = "-"
	}

	stats.Rank = rank.ComputeRank(stats.CurrentAverageCoffees)

	return stats, nil
}

func GetTotalMonths(period models.Period) int {
	totalMonths := 1
	if period.TotalMonths != 0 {
		totalMonths = period.TotalMonths
	}
	return totalMonths
}
