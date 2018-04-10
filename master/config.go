package master

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	. "github.com/Sypheos/fogflow/common/datamodel"
)

var (
	INFO     *log.Logger
	PROTOCOL *log.Logger
	ERROR    *log.Logger
	DEBUG    *log.Logger
)

//The default output for all the loggers is set to ioutil.Discard
func init() {
	INFO = log.New(ioutil.Discard, "", 0)
	PROTOCOL = log.New(ioutil.Discard, "", 0)
	ERROR = log.New(ioutil.Discard, "", 0)
	DEBUG = log.New(ioutil.Discard, "", 0)
}

type RepositoryCfg struct {
	URL     string `json:"url"`
	Backend string `json:"backend"`
}

type Config struct {
	MessageBus      string           `json:"message_bus"`
	MyIP            string           `json:"my_ip"`
	AgentPort       int              `json:"ngsi_agent_port"`
	IoTDiscoveryURL string           `json:"iot_discovery_url"`
	PLocation       PhysicalLocation `json:"physical_location"`
	Logging         struct {
		Info     string `json:"info"`
		Protocol string `json:"protocol"`
		Errlog   string `json:"error"`
		Debug    string `json:"debug"`
	} `json:"logging"`
}

var logTargets map[string]io.Writer = map[string]io.Writer{
	"stdout":  os.Stdout,
	"stderr":  os.Stderr,
	"discard": ioutil.Discard,
}

func (c *Config) SetLogTargets() {
	target, ok := logTargets[c.Logging.Info]
	if !ok {
		target = os.Stdout
	}
	INFO = log.New(target, "INFO: ", log.Ldate|log.Ltime)
	target, ok = logTargets[c.Logging.Protocol]
	if !ok {
		target = ioutil.Discard
	}
	PROTOCOL = log.New(target, "PROTOCOL: ", log.Ldate|log.Ltime)
	target, ok = logTargets[c.Logging.Errlog]
	if !ok {
		target = os.Stderr
	}
	ERROR = log.New(target, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	target, ok = logTargets[c.Logging.Debug]
	if !ok {
		target = ioutil.Discard
	}
	DEBUG = log.New(target, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func ParseConfig(confFile string, confVar *Config) error {
	file, err := os.Open(confFile)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(confVar)
	if err != nil {
		return err
	}

	return nil
}

func CreateConfig(configFile string) Config {
	var config Config

	if configFile == "" {
		config.MessageBus = "amqp://guest:guest@localhost:5672/"
	} else {
		err := ParseConfig(configFile, &config)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		}
	}

	config.SetLogTargets()

	return config
}
