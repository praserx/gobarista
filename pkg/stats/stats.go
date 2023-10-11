package stats

import (
	"errors"
	"fmt"

	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/models"
	"github.com/praserx/gobarista/pkg/stats/rank"
	"golang.org/x/exp/slices"
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
		TotalMonths:    currentPeriod.TotalMonths,
		TotalQuantity:  currentPeriod.TotalQuantity,
		TotalAverage:   currentPeriod.TotalQuantity / currentPeriod.TotalMonths / totalCustomers,
	}
}

func GetUserStats(uid uint, currentPeriod models.Period, currentBill models.Bill) (stats UserStats, err error) {
	prevBills, err := database.SelectAllBillsForUserExceptSpecifiedPeriod(uid, currentPeriod.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return UserStats{}, fmt.Errorf("billing period not found")
	} else if err != nil {
		return UserStats{}, fmt.Errorf("cannot get billing period: %v", err)
	}

	var issuedPeriods []uint
	for _, prevBill := range prevBills {
		stats.PrevTotalCoffees += prevBill.Quantity
		if !slices.Contains(issuedPeriods, prevBill.PeriodID) && prevBill.PeriodID != currentPeriod.ID {
			issuedPeriods = append(issuedPeriods, prevBill.PeriodID)
		}
	}

	var prevPeriods []models.Period
	for _, issuedPeriod := range issuedPeriods {
		prevPeriod, err := database.SelectPeriodByID(issuedPeriod)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UserStats{}, fmt.Errorf("billing period not found")
		} else if err != nil {
			return UserStats{}, fmt.Errorf("cannot get billing period: %v", err)
		}
		prevPeriods = append(prevPeriods, prevPeriod)
	}

	for _, prevPeriod := range prevPeriods {
		stats.PrevTotalMonths += prevPeriod.TotalMonths
	}

	if stats.PrevTotalMonths == 0 {
		stats.PrevAverageCoffees = 0
	} else {
		stats.PrevAverageCoffees = stats.PrevTotalCoffees / stats.PrevTotalMonths
	}

	stats.CurrentAverageCoffees = (stats.PrevTotalCoffees + currentBill.Quantity) / (stats.PrevTotalMonths + currentPeriod.TotalMonths)

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
