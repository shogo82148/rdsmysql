package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/shogo82148/rdsmysql/internal/certificate"
)

// Generate generates the configuration file for mysql.
func Generate(session *session.Session, dir string, config *Config) error {
	credentials := session.Config.Credentials
	token, err := rdsutils.BuildAuthToken(
		fmt.Sprintf("%s:%d", config.Host, config.Port),
		aws.StringValue(session.Config.Region),
		config.User,
		credentials,
	)
	if err != nil {
		return fmt.Errorf("fail to build auth token: %w", err)
	}
	pempath := filepath.Join(dir, "rds-combined-ca-bundle.pem")
	confpath := filepath.Join(dir, "my.conf")
	conf := fmt.Sprintf(`[client]
host = %s
user = %s
port = %d
password = %s
ssl-ca = %s
enable-cleartext-plugin
`, config.Host, config.User, config.Port, token, pempath)

	now := time.Now()
	if err := ioutil.WriteFile(fmt.Sprintf("%s.%d", confpath, now.UnixNano()), []byte(conf), 0600); err != nil {
		return err
	}
	if err := ioutil.WriteFile(fmt.Sprintf("%s.%d", pempath, now.UnixNano()), []byte(certificate.Certificate), 0600); err != nil {
		return err
	}

	if err := os.Rename(fmt.Sprintf("%s.%d", confpath, now.UnixNano()), confpath); err != nil {
		return err
	}
	if err := os.Rename(fmt.Sprintf("%s.%d", pempath, now.UnixNano()), pempath); err != nil {
		return err
	}

	return nil
}

// Config is the configuration for connecting to mysql servers.
type Config struct {
	User string
	Host string
	Port int
	Args []string
}

// Parse parses the args of mysql command.
func Parse(args []string) (*Config, error) {
	var conf Config
	conf.Port = 3306 // default port
	conf.Args = make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-u", "--user":
			i++
			if i >= len(args) {
				return nil, errors.New("invalid user option")
			}
			conf.User = args[i]
		case "-h", "--host":
			i++
			if i >= len(args) {
				return nil, errors.New("invalid host option")
			}
			conf.Host = args[i]
		case "-P", "--port":
			i++
			if i >= len(args) {
				return nil, errors.New("invalid port option")
			}
			port, err := strconv.Atoi(args[i])
			if err != nil {
				return nil, errors.New("fail to parse port")
			}
			conf.Port = port
		default:
			conf.Args = append(conf.Args, args[i])
		}
	}
	return &conf, nil
}
