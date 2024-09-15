package usecase

import (
	"errors"
	"fmt"
	. "github.com/craftizmv/rewards/internal/app/repository"
	"github.com/craftizmv/rewards/internal/app/usecase/helper"
	"github.com/craftizmv/rewards/internal/data/dtos"
	. "github.com/craftizmv/rewards/internal/data/infrastructure/cache"
	"github.com/craftizmv/rewards/internal/data/infrastructure/external/proxies"
	"github.com/craftizmv/rewards/internal/data/infrastructure/queue/events"
	"github.com/craftizmv/rewards/internal/domain/entities"
	"github.com/craftizmv/rewards/pkg/logger"
)

// RewardUseCaseImpl implements the gift-related use cases
type RewardUseCaseImpl struct {
	cache      ICache[entities.Order]
	rewardRepo RewardRepository
	log        logger.ILogger
	proxies    *RewardProxies
}

type RewardProxies struct {
	CampaignProxy  *proxies.CampaignProxy
	InventoryProxy *proxies.InventoryProxy
	EmailProxy     *proxies.EmailProxy
	ShippingProxy  *proxies.ShippingProxy
	UserProxy      *proxies.UserProxy
	OrderProxy     *proxies.OrderProxy
}

// NewRewardUseCaseImpl injects dependencies into the RewardUseCaseImpl
func NewRewardUseCaseImpl(cache ICache[entities.Order], rewardRepo RewardRepository, log logger.ILogger, proxies *RewardProxies) *RewardUseCaseImpl {
	return &RewardUseCaseImpl{
		cache:      cache,
		rewardRepo: rewardRepo,
		log:        log,
		proxies:    proxies,
	}
}

// AllocateReward allocateGift allocates a gift based on the order ID
func (rewardUseCase *RewardUseCaseImpl) AllocateReward(event events.AllocateReward) error {
	// retrieve order info from the shared cache.
	order, found := rewardUseCase.cache.Get(helper.GetOrderKey(event.OrderID))
	if !found {
		// TODO : If not found in cache .. may be we can check in DB or simply reject allocating, trigger a background update if cache.
		rewardUseCase.log.Error("order not found", "orderID", event.OrderID)
		return errors.New("order not found")
	}

	// checking from the order object if the reward is already issued
	if order.RewardStatus != entities.RewardStatusNone {
		return errors.New("reward is already processed")
	}

	if order.IsOrderEffectivelyRolledBack() {
		return errors.New("order is already rolled back")
	}

	// Get the list of productIDs for the rewardGroup
	productIDList, err := rewardUseCase.rewardRepo.GetProductIDsFromRewardGroup(event.RewardTypeID)
	if err != nil {
		return err
	}

	// Assumption : Inventory Proxy provides an API to check the inventory
	// availability of items needed to be allocated as part of the reward.
	ok, _ := rewardUseCase.proxies.InventoryProxy.BulkVerifyInventoryAvailability(productIDList)
	if !ok {
		return errors.New("can not allocate reward, inventory unavailable")
	}

	// 2. Block inventory and update the order cache.
	allOK, itemIDList := rewardUseCase.proxies.InventoryProxy.BlockInventoryForProducts(productIDList)
	if !allOK {
		return errors.New("could not block inventory")
	}

	// congrats : all check passed, create mappings for the rewardGroup.
	// insert to reward group reward item mapping
	rewardGroupID := helper.GenerateRandomInt64() // TODO : This needs to be updated in the order table.
	err = rewardUseCase.rewardRepo.InsertRewardGroupRewardItemsBatch(rewardGroupID, itemIDList, 5)
	if err != nil {
		rewardUseCase.log.Error("failed to insert in reward group reward item table", "error", err)
		return err
	}

	//insert to order_reward_item mapping.
	// NOTE : Update the shipment async when the reward is allocated - we can retry and keep retrying until it is success. (also, issue alert)
	userDetail := rewardUseCase.proxies.UserProxy.GetUserDetails(event.UserID)
	shipmentResponse, err := rewardUseCase.proxies.ShippingProxy.ShipItems(itemIDList, userDetail)
	if err != nil {
		// TODO: Handle various kinds of error
		rewardUseCase.log.Error("failed to ship items", "error", err)
		return err
	}

	// using some random num.
	if shipmentResponse.Cost > 100000 {
		rewardUseCase.log.Error("failed to ship items, cost too high", "cost", shipmentResponse.Cost)
		return errors.New("failed to ship items, cost too high")
	}

	err = rewardUseCase.rewardRepo.UpdateOrderRewardItemsBatch(helper.CreateOrderRewardItems(event.OrderID, itemIDList), 5)
	if err != nil {
		rewardUseCase.log.Error("failed to update order reward item table", "error", err)
		return err
	}

	// TODO : Order cache . - use order proxy to do that.

	// TODO: Send email
	err = rewardUseCase.proxies.EmailProxy.SendEmail(userDetail.UserName, userDetail.Email, "Hi, XZY")
	if err != nil {
		rewardUseCase.log.Error("failed to send email", "error", err)
		// see if we handle retry or communicate via whatsapp etc.
	}

	rewardUseCase.log.Info("successfully allocated reward", "orderID", event.OrderID)

	return nil
}

// CancelReward cancels the reward associated with the order ID
func (rewardUseCase *RewardUseCaseImpl) CancelReward(revokeReward events.RevokeReward) error {

	// TODO : Check the shipping status of the Reward, If valid then proceed further.

	// TODO: Make all below steps transaction to avoid data inconsistency.
	// Get RewardGroupID
	rewardGroupIDs, err := rewardUseCase.rewardRepo.GetRewardGroupIDByOrderID(revokeReward.OrderID)
	if err != nil {
		rewardUseCase.log.Error("failed to find reward group", "error", err)
		return err
	}

	if len(rewardGroupIDs) <= 0 {
		rewardUseCase.log.Error("failed to find reward group", "error", err)
		return errors.New("failed to find reward group")
	}

	// update order cache
	_, err = rewardUseCase.proxies.OrderProxy.UpdateOrderRewardStatus(revokeReward.OrderID, rewardGroupIDs[0], string(entities.RewardStatusCancelled))
	if err != nil {
		rewardUseCase.log.Error("failed to update order reward status", "error", err)
		return err
	}

	// removing the relationship of reward with order
	err = rewardUseCase.rewardRepo.DeleteRewardGroupByOrderID(revokeReward.OrderID, rewardGroupIDs[0])
	if err != nil {
		rewardUseCase.log.Error("failed to delete reward group", "error", err)
		return err
	}

	err = rewardUseCase.rewardRepo.DeleteRewardItemsByOrderID(revokeReward.OrderID)
	if err != nil {
		rewardUseCase.log.Error("failed to delete reward items", "error", err)
		return err
	}

	//NOTE : Don't delete the generated reward, as this can be used for re-allocation. Can be cleanup later by a JOB.
	return nil
}

// ReAllocateReward reAllocateGift reallocates a gift for the given order ID
func (rewardUseCase *RewardUseCaseImpl) ReAllocateReward(reAllocateEvent events.ReAllocateReward) error {
	// TODO: Logic Same as Allocate Reward from the buffer queue to pick the new order and execute set of steps to
	// TODO: update the entries in transactional mode

	return nil
}

func (rewardUseCase *RewardUseCaseImpl) CheckRewardEligibility(orderDTO *dtos.OrderDTO) (bool, error) {
	// Here we are checking below 4 conditions:
	// 1. campaign is active
	// 2. order val satisfies the eligibility criteria
	// 3. RewardItem inventory availability
	// 4. correct order status - cache.

	// 1. campaign is active
	// TODO : check for proxies null condition if needed.
	campaign := rewardUseCase.proxies.CampaignProxy.FetchMostEligibleCampaign()
	if campaign == nil {
		return false, errors.New("failed to fetch campaign")
	}

	if campaign.Status != dtos.Active {
		return false, errors.New("campaign is not active")
	}

	if campaign.AllocatedRewards == campaign.TotalEligibleRewards {
		return false, errors.New("rewardGroup allocation limit exhausted")
	}

	if orderDTO.OrderValue < campaign.EligibilityCriteria.MinimumPurchaseAmount {
		return false, errors.New("order value is too small")
	}

	// check inventory of the rewardGroup from the inventory table.
	// Check the availability of productID obtained from RewardGroup data.
	if rewardGroup, err := rewardUseCase.rewardRepo.GetRewardGroupByID(campaign.RewardGroupID); err != nil {
		return false, fmt.Errorf("failed to get rewardGroup from repository: %w", err)
	} else {
		if rewardGroup == nil {
			return false, fmt.Errorf("no rewardGroup found for rewardGroup ID %d", campaign.RewardGroupID)
		} else {
			// loop over rewardItems in rewardGroup and check the inventory availability
			// TODO Get productList from RewardGroup, then using the product list check inventory (as done in AllocateReward func)
		}
	}

	if orderCacheObj, ok := rewardUseCase.cache.Get(helper.GetOrderKey(orderDTO.OrderID)); !ok {
		return false, fmt.Errorf("failed to get order from cach")
	} else {
		if orderCacheObj.IsComplete() {
			return false, errors.New("order is already completed")
		}
	}

	return true, nil
}
