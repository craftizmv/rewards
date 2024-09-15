package dtos

// RewardEligibilityResponse represents the response of reward eligibility check.
type RewardEligibilityResponse struct {
	Eligible bool   `json:"eligible"` // Eligibility status (true or false)
	Message  string `json:"message"`  // A message providing more context
	Reason   string `json:"reason"`   // (Optional) Reason for ineligibility
}
