package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"gopkg.in/yaml.v2"
	"os"
	"provider_mq/internal/consts"
	"strconv"
	"time"
)

type Config struct {
	QueueName  string     `yaml:"queueName"`
	ReplyTo    string     `yaml:"replyTo"`
	Topic      string     `yaml:"topic" default:""`
	BindingKey string     `yaml:"bindingKey" default:""`
	DLE        bool       `yaml:"dle" default:"false"`
	ArgsQueue  amqp.Table `yaml:"argQueue"`
}

func init() {
	viper.SetConfigFile(consts.EnvFile)
	viper.AddConfigPath(consts.EnvFileDirectory)
	err := viper.ReadInConfig()
	if err != nil {
		log.Debug().Err(err).
			Msg("Error occurred while reading env file, might fallback to OS env config")
	}
	viper.AutomaticEnv()
}

// GetQueueConf function can be used to get configs value for PublishConnection
// Modify ./queue_config.yaml (consts.QueuesConf) for change params
func GetQueueConf() []Config {
	var configs []Config

	source, err := os.ReadFile(consts.QueuesConf)

	if err != nil {
		log.Debug().Msgf("failed reading config file: %v\n", err)

		configs = append(configs, Config{
			QueueName:  "fib",
			ReplyTo:    "response",
			Topic:      "",
			BindingKey: ""})

	} else {
		err = yaml.Unmarshal(source, &configs)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}

	if len(configs) > 10 {
		log.Error().Msgf("Please! Do not exceed 10 queues...")
		configs = configs[:len(configs)-1]
	}

	fmt.Printf("config:\n%+v\n", configs)

	return configs
}

// GetEnvVar This function can be used to get ENV Var in our App
// Modify this if you change the library to read ENV
func GetEnvVar(name string) string {
	if !viper.IsSet(name) {
		log.Debug().Msgf("Environment variable %s is not set", name)
		return ""
	}
	value := viper.GetString(name)
	return value
}

func GetTlsConf() *tls.Config {
	caCert, err := os.ReadFile(GetEnvVar("_CACERT"))
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read CACert")
	}

	cert, err := tls.LoadX509KeyPair(GetEnvVar("_CERT"), GetEnvVar("_KEY"))
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read Certificate, Key")
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(caCert)

	tlsConf := &tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{cert},
	}

	return tlsConf
}

func GetCorrelationId() string {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	return "msg" + strconv.FormatInt(t, 10)
}

func Timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
