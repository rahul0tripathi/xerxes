package datastore

func NotifyChange(channel string, message string) {
	RedisPubSub.Publish(PubSubContext, channel, message)
}
