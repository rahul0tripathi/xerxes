package datastore

//var (
//	RedisClient *redisv8.Client
//	RedisPubSub *redisv8.Client
//	PubSubContext = context.Background()
//)
//
//// Initializes redis connection
//func initClient() {
//	RedisClient = redisv8.NewClient(&redisv8.Options{
//		Addr:     config.CacheConfig.Addr,
//		Password: config.CacheConfig.Password,
//		DB:       config.CacheConfig.DB,
//	})
//	RedisPubSub = redisv8.NewClient(&redisv8.Options{
//		Addr:     config.CacheConfig.Addr,
//		Password: config.CacheConfig.Password,
//		DB:       config.CacheConfig.DB,
//	})
//
//}
//
//func init() {
//	err := config.LoadConfig.Cache()
//	if err != nil {
//		panic(err)
//	}
//	initClient()
//}
