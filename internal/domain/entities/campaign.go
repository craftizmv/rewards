package entities

import (
	"errors"
	"fmt"
	. "github.com/craftizmv/rewards/internal/data/dtos"
	"github.com/google/uuid"
	"time"
)

type Campaign struct {
	ID                   uuid.UUID
	Name                 string
	StartDate            time.Time
	EndDate              time.Time
	Status               Status // maybe we can map the status to a domain based status
	Budget               float64
	AllocatedRewards     int
	TotalEligibleRewards int
	TargetAudience       string
}

// Business invariants
var (
	ErrInvalidDateRange        = errors.New("start date must be before end date")
	ErrInvalidStatus           = errors.New("invalid campaign status")
	ErrInvalidBudget           = errors.New("budget must be greater than zero")
	ErrGiftExceedsBudget       = errors.New("allocated gifts exceed budget")
	ErrNegativeMaxGiftsPerUser = errors.New("max gifts per user cannot be negative")
	ErrInvalidActivation       = errors.New("campaign cannot be activated outside of the date range")
)

// CampaignValidator defines the interface for validating CampaignDTOs.
type CampaignValidator interface {
	Validate(campaign *CampaignDTO, currentTime time.Time) error
}

// campaignValidator is the concrete implementation of CampaignValidator.
type campaignValidator struct{}

// NewCampaignValidator creates a new instance of campaignValidator.
func NewCampaignValidator() CampaignValidator {
	return &campaignValidator{}
}

// Validate checks if the provided CampaignDTO adheres to all business invariants.
func (v *campaignValidator) Validate(campaign *CampaignDTO, currentTime time.Time) error {
	if campaign.StartDate.After(campaign.EndDate) {
		return fmt.Errorf("%w: start date (%s) is after end date (%s)", ErrInvalidDateRange, campaign.StartDate, campaign.EndDate)
	}

	switch campaign.Status {
	case Active, Paused, Ended:
		// Valid statuses
	default:
		return fmt.Errorf("%w: received status '%s'", ErrInvalidStatus, campaign.Status)
	}

	if campaign.Budget <= 0 {
		return fmt.Errorf("%w: budget (%f) is not greater than zero", ErrInvalidBudget, campaign.Budget)
	}

	// Assuming each gift has a fixed value, e.g., $10. Modify as per actual business logic.
	const giftValue = 10.0
	totalGiftValue := float64(campaign.AllocatedRewards) * giftValue
	if totalGiftValue > campaign.Budget {
		return fmt.Errorf("%w: total gift value (%f) exceeds budget (%f)", ErrGiftExceedsBudget, totalGiftValue, campaign.Budget)
	}

	if campaign.MaxGiftsPerUser < 0 {
		return fmt.Errorf("%w: max gifts per user (%d) is negative", ErrNegativeMaxGiftsPerUser, campaign.MaxGiftsPerUser)
	}

	// If the campaign is active, ensure current time is within the campaign's date range.
	if campaign.Status == Active {
		if currentTime.Before(campaign.StartDate) || currentTime.After(campaign.EndDate) {
			return fmt.Errorf("%w: current time (%s) is outside the campaign date range (%s - %s)", ErrInvalidActivation, currentTime, campaign.StartDate, campaign.EndDate)
		}
	}

	return nil
}
