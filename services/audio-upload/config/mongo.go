package config

type MongoConfig struct {
	URI                string `env:"MONGO_URI" env-default:"mongodb://localhost:27017"`
	Database           string `env:"MONGO_DATABASE" env-default:"luna_audio"`
	TracksCollection   string `env:"MONGO_TRACKS_COLLECTION" env-default:"tracks"`
}