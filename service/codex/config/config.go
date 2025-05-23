package config

import (
	"github.com/charmingruby/doris/lib/config"
)

type Config config.Config[CustomConfig]

type CustomConfig struct {
	RestServerHost               string `env:"REST_SERVER_HOST" envDefault:"localhost"`
	RestServerPort               string `env:"REST_SERVER_PORT" envDefault:"3000"`
	AWSRegion                    string `env:"AWS_REGION,required"`
	AWSEmbeddingSourceDocsBucket string `env:"AWS_EMBEDDING_SOURCE_DOCS_BUCKET,required"`
	DatabaseHost                 string `env:"DATABASE_HOST,required"`
	DatabasePort                 string `env:"DATABASE_PORT,required"`
	DatabaseUser                 string `env:"DATABASE_USER,required"`
	DatabasePassword             string `env:"DATABASE_PASSWORD,required"`
	DatabaseName                 string `env:"DATABASE_NAME,required"`
	DatabaseSSL                  string `env:"DATABASE_SSL,required"`
	NatsStream                   string `env:"NATS_STREAM"`
	APIKeyActivatedTopic         string `env:"API_KEY_ACTIVATED_TOPIC"`
	APIKeyDelegatedTopic         string `env:"API_KEY_DELEGATED_TOPIC"`
	CodexDocumentUploadedTopic   string `env:"CODEX_DOCUMENT_UPLOADED_TOPIC"`
	JWTSecret                    string `env:"JWT_SECRET,required"`
	JWTIssuer                    string `env:"JWT_ISSUER,required"`
	OllamaEmbeddingModel         string `env:"OLLAMA_EMBEDDING_MODEL,required"`
	OllamaCompletionModel        string `env:"OLLAMA_COMPLETION_MODEL,required"`
	OllamaBaseURL                string `env:"OLLAMA_BASE_URL,required"`
}

func New() (Config, error) {
	cfg, err := config.New[CustomConfig]()
	if err != nil {
		return Config{}, err
	}

	return Config(cfg), nil
}
