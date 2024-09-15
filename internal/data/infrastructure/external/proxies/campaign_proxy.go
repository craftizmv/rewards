package proxies

import (
	. "github.com/craftizmv/rewards/internal/data/dtos"
	"github.com/craftizmv/rewards/internal/data/infrastructure/external/mocks"
)

type CampaignProxy struct{}

func NewCampaignProxy() *CampaignProxy {
	return &CampaignProxy{}
}

// ReturnMockCampaigns - returns list of mock campaigns
func (p *CampaignProxy) ReturnMockCampaigns() []*CampaignDTO {
	return mocks.MockCampaigns()
}

func (p *CampaignProxy) FetchMostEligibleCampaign() *CampaignDTO {
	return mocks.MockValidCampaign()
}
