package config

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		in                []string
		user              string
		host              string
		port              int
		useSystemCertPool bool
		out               []string
	}{
		{
			in:   []string{},
			user: "",
			host: "",
			port: 3306,
			out:  []string{},
		},

		// short option
		{
			in:   []string{"-u", "username", "-h", "rds.example.com", "-P", "12345"},
			user: "username",
			host: "rds.example.com",
			port: 12345,
			out:  []string{},
		},

		// long option
		{
			in:   []string{"--user", "username", "--host", "rds.example.com", "--port", "12345"},
			user: "username",
			host: "rds.example.com",
			port: 12345,
			out:  []string{},
		},

		{
			in:   []string{"-u", "username", "-abc", "--foobar", "foobar"},
			user: "username",
			host: "",
			port: 3306,
			out:  []string{"-abc", "--foobar", "foobar"},
		},

		{
			in:                []string{"-u", "username", "--use-system-cert-pool"},
			user:              "username",
			host:              "",
			port:              3306,
			useSystemCertPool: true,
			out:               []string{},
		},
	}

	for i, tc := range cases {
		conf, err := Parse(tc.in)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}
		if conf.User != tc.user {
			t.Errorf("%d: want %s, got %s", i, tc.user, conf.User)
		}
		if conf.Host != tc.host {
			t.Errorf("%d: want %s, got %s", i, tc.host, conf.Host)
		}
		if conf.Port != tc.port {
			t.Errorf("%d: want %d, got %d", i, tc.port, conf.Port)
		}
		if conf.UseSystemCertPool != tc.useSystemCertPool {
			t.Errorf("%d: want %v, got %v", i, tc.useSystemCertPool, conf.UseSystemCertPool)
		}
		if !reflect.DeepEqual(conf.Args, tc.out) {
			t.Errorf("%d: want %#v, got %#v", i, tc.out, conf.Args)
		}
	}
}
