package constants

type EnvKeys struct {
    Env           string
    ServerAddress string

    // Mongo
    MongoHost string
    MongoPort string
    MongoUser string
    MongoPass string
    MongoDB   string
}

var Env = EnvKeys{
    Env:           "ENV",
    ServerAddress: "SERVER_ADDRESS",

    MongoHost: "MONGO_HOST",
    MongoPort: "MONGO_PORT",
    MongoUser: "MONGO_USER",
    MongoPass: "MONGO_PASS",
    MongoDB:   "MONGO_DB_NAME",
}
