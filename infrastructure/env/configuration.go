package env

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var (
	once   sync.Once
	config = &configuration{}
)

type configuration struct {
	App   App   `json:"app"`
	DB    DB    `json:"db"`
	Files Files `json:"files"`
	Aws   Aws   `json:"aws"`
}

type App struct {
	ServiceName       string `json:"service_name"`
	PathLog           string `json:"path_log"`
	LogReviewInterval int    `json:"log_review_interval"`
	RegisterLog       bool   `json:"register_log"`
	WorkerInterval    int    `json:"worker_interval"`
}

type DB struct {
	Engine   string `json:"engine"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Instance string `json:"instance"`
	IsSecure bool   `json:"is_secure"`
	SSLMode  string `json:"ssl_mode"`
}

type Files struct {
	Repo string `json:"repo"`
	S3   struct {
		Bucket string `json:"bucket"`
		Region string `json:"region"`
	} `json:"s3"`
}

type Aws struct {
	AWSACCESSKEYID     string `json:"AWS_ACCESS_KEY_ID"`
	AWSSECRETACCESSKEY string `json:"AWS_SECRET_ACCESS_KEY"`
	AWSDEFAULTREGION   string `json:"AWS_DEFAULT_REGION"`
}

func NewConfiguration() *configuration {
	fromFile()
	return config
}

// LoadConfiguration lee el archivo configuration.json
// y lo carga en un objeto de la estructura Configuration
func fromFile() {
	once.Do(func() {
		b, err := os.ReadFile("config.json")
		if err != nil {
			log.Fatalf("no se pudo leer el archivo de configuración: %s", err.Error())
		}

		err = json.Unmarshal(b, config)
		if err != nil {
			log.Fatalf("no se pudo parsear el archivo de configuración: %s", err.Error())
		}
	})
}
