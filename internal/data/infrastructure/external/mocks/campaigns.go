package mocks

import (
	. "github.com/craftizmv/rewards/internal/data/dtos"
	uuid "github.com/satori/go.uuid"
	"time"
)

// MockCampaigns returns a list of mock CampaignDTOs for testing purposes.
func MockCampaigns() []*CampaignDTO {
	return []*CampaignDTO{
		{
			ID:                   uuid.NewV4(),
			Name:                 "Black Friday Sale",
			StartDate:            time.Now().Add(-time.Hour * 24 * 10), // Started 10 days ago
			EndDate:              time.Now().Add(time.Hour * 24 * 5),   // Ends in 5 days
			Status:               Active,
			Budget:               10000.0,
			AllocatedRewards:     100,
			TargetAudience:       "All Customers",
			TotalEligibleRewards: 1000,
		},
		{
			ID:                   uuid.NewV4(),
			Name:                 "Holiday Discounts",
			StartDate:            time.Now().Add(time.Hour * 24 * 30), // Starts in 30 days
			EndDate:              time.Now().Add(time.Hour * 24 * 60), // Ends in 60 days
			Status:               Paused,                              // Paused until it starts
			Budget:               5000.0,
			AllocatedRewards:     50,
			TargetAudience:       "Loyal Customers",
			TotalEligibleRewards: 1000,
		},
		{
			ID:                   uuid.NewV4(),
			Name:                 "Summer Clearance",
			StartDate:            time.Now().Add(-time.Hour * 24 * 60), // Started 60 days ago
			EndDate:              time.Now().Add(-time.Hour * 24 * 10), // Ended 10 days ago
			Status:               Ended,
			Budget:               8000.0,
			AllocatedRewards:     200,
			TargetAudience:       "Young Adults",
			TotalEligibleRewards: 1000,
		},
		{
			ID:                   uuid.NewV4(),
			Name:                 "Invalid Date Campaign",
			StartDate:            time.Now().Add(time.Hour * 24 * 10), // Starts in 10 days
			EndDate:              time.Now().Add(time.Hour * 24 * 5),  // Ends before it starts
			Status:               Paused,
			Budget:               3000.0,
			AllocatedRewards:     20,
			TargetAudience:       "New Customers",
			TotalEligibleRewards: 1000,
		},
		{
			ID:                   uuid.NewV4(),
			Name:                 "Gift Exceeds Budget",
			StartDate:            time.Now().Add(-time.Hour * 24 * 1), // Started 1 day ago
			EndDate:              time.Now().Add(time.Hour * 24 * 15), // Ends in 15 days
			Status:               Active,
			Budget:               100.0, // Budget is too small
			AllocatedRewards:     50,    // Gifts exceed the budget
			TargetAudience:       "Students",
			TotalEligibleRewards: 1000,
		},
		{
			ID:                   uuid.NewV4(),
			Name:                 "Zero Budget Campaign",
			StartDate:            time.Now().Add(-time.Hour * 24 * 1), // Started 1 day ago
			EndDate:              time.Now().Add(time.Hour * 24 * 10), // Ends in 10 days
			Status:               Active,
			Budget:               0.0, // Zero budget
			AllocatedRewards:     0,
			TargetAudience:       "General Audience",
			TotalEligibleRewards: 1000,
		},
	}
}

// MockValidCampaign returns a single valid mock
func MockValidCampaign() *CampaignDTO {
	return &CampaignDTO{
		ID:                   uuid.NewV4(),
		Name:                 "New Year Promotion",
		StartDate:            time.Now().Add(-time.Hour * 24 * 2), // Started 2 days ago
		EndDate:              time.Now().Add(time.Hour * 24 * 7),  // Ends in 7 days
		Status:               Active,
		Budget:               15000.0,
		AllocatedRewards:     75,
		TargetAudience:       "All Customers",
		TotalEligibleRewards: 1000,
	}
}

// MockInvalidCampaign returns a single invalid mock campaign for testing purposes.
func MockInvalidCampaign() *CampaignDTO {
	return &CampaignDTO{
		ID:                   uuid.NewV4(),
		Name:                 "Broken Campaign",
		StartDate:            time.Now().Add(time.Hour * 24 * 10), // Starts in 10 days
		EndDate:              time.Now().Add(time.Hour * 24 * 5),  // Ends before it starts
		Status:               Paused,
		Budget:               1000.0,
		AllocatedRewards:     10,
		TargetAudience:       "VIP Customers",
		TotalEligibleRewards: 1000,
	}
}
