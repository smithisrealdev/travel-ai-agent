package agents

import "log"

// BudgetPlan represents the breakdown of a travel budget
type BudgetPlan struct {
	Flight    int `json:"flight"`
	Hotel     int `json:"hotel"`
	Food      int `json:"food"`
	Transport int `json:"transport"`
	Misc      int `json:"misc"`
}

// EstimateBudget splits a total budget into travel expense categories
// Parameters:
//   - total: Total budget amount in THB
// Returns:
//   - BudgetPlan with calculated amounts for each category
func EstimateBudget(total int) BudgetPlan {
	// Handle edge case: negative or zero budget
	if total <= 0 {
		log.Printf("Warning: Invalid total budget %d, returning zero budget", total)
		return BudgetPlan{
			Flight:    0,
			Hotel:     0,
			Food:      0,
			Transport: 0,
			Misc:      0,
		}
	}

	// Calculate each category based on percentages
	flight := int(float64(total) * 0.45)    // 45%
	hotel := int(float64(total) * 0.25)     // 25%
	food := int(float64(total) * 0.15)      // 15%
	transport := int(float64(total) * 0.10) // 10%
	misc := int(float64(total) * 0.05)      // 5%

	budget := BudgetPlan{
		Flight:    flight,
		Hotel:     hotel,
		Food:      food,
		Transport: transport,
		Misc:      misc,
	}

	log.Printf("Budget estimated for total %d THB: Flight=%d, Hotel=%d, Food=%d, Transport=%d, Misc=%d",
		total, budget.Flight, budget.Hotel, budget.Food, budget.Transport, budget.Misc)

	return budget
}
