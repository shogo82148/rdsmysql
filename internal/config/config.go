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
		// root CA. The mysql client, unlike Go's crypto/x509, has no notion of "the system
		// cert pool" and refuses to verify a certificate at all (ssl-mode=VERIFY_IDENTITY)
		// without an explicit ssl-ca file, so one of the OS's own CA bundle files is located
		// and passed as ssl-ca.
		caPath, err := findSystemCertBundle()
		if err != nil {
			return err
		}
		conf = fmt.Sprintf(`[client]
host = %s
user = %s
port = %d
password = %s
ssl-ca = %s
ssl-mode = VERIFY_IDENTITY
enable-cleartext-plugin
`, config.Host, config.User, config.Port, token, caPath)
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

// systemCertBundlePaths lists locations of CA certificate bundle files maintained by the
// operating system or its package manager, checked in order. It mirrors the candidates the Go
// standard library checks for x509.SystemCertPool on Unix (see crypto/x509/root_linux.go),
// which aren't exported for reuse here, plus the path macOS itself maintains.
var systemCertBundlePaths = []string{
	"/etc/ssl/cert.pem",                                 // macOS, Alpine Linux
	"/etc/ssl/certs/ca-certificates.crt",                // Debian, Ubuntu, Gentoo, Arch
	"/etc/pki/tls/certs/ca-bundle.crt",                  // Fedora, RHEL 6
	"/etc/ssl/ca-bundle.pem",                            // OpenSUSE
	"/etc/pki/tls/cacert.pem",                           // OpenELEC
	"/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem", // CentOS, RHEL 7+, Amazon Linux 2+
}

// findSystemCertBundle returns the path of the first existing CA certificate bundle file from
// systemCertBundlePaths.
func findSystemCertBundle() (string, error) {
	for _, path := range systemCertBundlePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", errors.New("could not find the system's CA certificate bundle for --use-system-cert-pool")
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
