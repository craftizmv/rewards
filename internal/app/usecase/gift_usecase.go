package usecase

type GiftUseCases interface {
	allocateGift() error
	cancelGift() error
	reAllocateGift() error
}
