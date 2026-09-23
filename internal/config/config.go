package config

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/shogo82148/rdsmysql/v2"
)

// Generate generates the configuration file for mysql.
func Generate(ctx context.Context, awsConfig aws.Config, dir string, config *Config) error {
	cred := awsConfig.Credentials
	region := awsConfig.Region
	if region == "" {
		return errors.New("region is not specified")
	}
	endpoint := net.JoinHostPort(config.Host, strconv.Itoa(config.Port))
	token, err := auth.BuildAuthToken(ctx, endpoint, region, config.User, cred)
	if err != nil {
		return fmt.Errorf("fail to build auth token: %w", err)
	}
	confpath := filepath.Join(dir, "my.conf")
	now := time.Now()

	var conf string
	if config.UseSystemCertPool {
		// Trust the system's CA certificates instead of the bundled Amazon RDS root
		// certificates. Enable this when connecting through Amazon RDS Proxy, which presents
		// certificates issued by AWS Certificate Manager (ACM) rather than the Amazon RDS
		// root CA. ssl-mode is set explicitly because without ssl-ca, mysql otherwise
		// defaults to ssl-mode=PREFERRED, which doesn't verify the server certificate at all.
		conf = fmt.Sprintf(`[client]
host = %s
user = %s
port = %d
password = %s
ssl-mode = VERIFY_IDENTITY
enable-cleartext-plugin
`, config.Host, config.User, config.Port, token)
	} else {
		pempath := filepath.Join(dir, "rds-combined-ca-bundle.pem")
		if err := os.WriteFile(fmt.Sprintf("%s.%d", pempath, now.UnixNano()), []byte(rdsmysql.Certificates), 0600); err != nil {
			return err
		}
		if err := os.Rename(fmt.Sprintf("%s.%d", pempath, now.UnixNano()), pempath); err != nil {
			return err
		}
		conf = fmt.Sprintf(`[client]
host = %s
user = %s
port = %d
password = %s
ssl-ca = %s
enable-cleartext-plugin
`, config.Host, config.User, config.Port, token, pempath)
	}

	if err := os.WriteFile(fmt.Sprintf("%s.%d", confpath, now.UnixNano()), []byte(conf), 0600); err != nil {
		return err
	}
	if err := os.Rename(fmt.Sprintf("%s.%d", confpath, now.UnixNano()), confpath); err != nil {
		return err
	}

	return nil
}

// Config is the configuration for connecting to mysql servers.
type Config struct {
	User    string
	Host    string
	Port    int
	Version bool

	// UseSystemCertPool makes the connection trust the system's CA certificate pool
	// instead of the Amazon RDS root certificates. Enable this when connecting through
	// Amazon RDS Proxy, which presents certificates issued by AWS Certificate Manager (ACM)
	// rather than the Amazon RDS root CA.
	UseSystemCertPool bool

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
		case "-V", "--version":
			conf.Version = true
		case "--use-system-cert-pool":
			conf.UseSystemCertPool = true
		default:
			conf.Args = append(conf.Args, args[i])
		}
	}
	return &conf, nil
}
