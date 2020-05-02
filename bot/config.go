package bot

const (
	// MongoConnectionTime the time in seconds that the
	// bot has to connect to the mongo database before
	// timeingout.
	MongoConnectionTime = 10

	// MongoQueryTimeout the time in seconds that the
	// mongo database has to perform a query before it
	// times out.
	MongoQueryTimeout = 2

	// MongoPingTime is the time the mongo database has
	// to respond to a ping before
	MongoPingTime = 2

	// PingInterval is the time between each ping to check
	// that the Bot is still connected to the mongodb.
	PingInterval = 20
)
