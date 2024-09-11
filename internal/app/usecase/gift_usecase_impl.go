package usecase

type GiftUseCaseImpl struct {
	// mention dependencies
}

// NewGiftUseCaseImpl - inject dependency in the constructor
func NewGiftUseCaseImpl() *GiftUseCaseImpl {
	return &GiftUseCaseImpl{}
}

func (i *GiftUseCaseImpl) allocateGift() error {
	return nil
}

func (i *GiftUseCaseImpl) cancelGift() error {
	return nil
}

func (i *GiftUseCaseImpl) reAllocateGift() error {
	return nil
}
