package publisher

//go:generate mockery --name IPublisher
type IPublisher interface {
	PublishMessage(msg interface{}) error
	IsPublished(msg interface{}) bool
}
