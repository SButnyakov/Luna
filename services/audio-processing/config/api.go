package config

type API struct {
	V1 struct {
		Root  string `env:"HTTP_API_V1" env-default:"/v1"`
		Audio struct {
			Root    string `env:"HTTP_API_AUDIO_GROUP" env-default:"/audio"`
			Process string `env:"HTTP_API_AUDIO_PROCESS" env-default:"/process"`
		}
	}
}
