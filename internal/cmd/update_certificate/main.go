package main

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"go/format"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := downloadCertificate(ctx, &options{
		file: "rds.go",
		pkg:  "certificate",
		url:  "https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem",
		name: "rdsCertificates",
		comment: `// rdsCertificates is the intermediate and root [certificates] for [Amazon RDS MySQL] and [Amazon Aurora MySQL].
//
// [Amazon RDS MySQL]: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.SSL.html#UsingWithRDS.SSL.CertificatesAllRegions
// [Amazon Aurora MySQL]: https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.SSL.html#UsingWithRDS.SSL.CertificatesAllRegions
// [certificates]: https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem
`,
	})
	if err != nil {
		panic(err)
	}
}

type options struct {
	file    string
	pkg     string
	url     string
	name    string
	comment string
}

func downloadCertificate(ctx context.Context, opts *options) error {
	pemCerts, err := download(ctx, opts.url)
	if err != nil {
		return err
	}

	certs, err := parseCertificate(pemCerts)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	buf.WriteString("package " + opts.pkg + "\n\n")
	buf.WriteString(opts.comment)
	buf.WriteString("const " + opts.name + " = `")
	buf.Write(pemCerts)
	buf.WriteString("`\n\n")

	buf.WriteString("// " + opts.name + " contains:\n")
	buf.WriteString("//\n")
	for _, cert := range certs {
		nbf := cert.NotBefore.Format(time.RFC3339)
		naf := cert.NotAfter.Format(time.RFC3339)
		fmt.Fprintf(buf, "// - %50s (not before: %s, not after: %s)\n", cert.Subject.CommonName, nbf, naf)
	}

	data, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	return os.WriteFile(opts.file, data, 0644)
}

func download(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func parseCertificate(pemCerts []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		certBytes := block.Bytes
		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}
	return certs, nil
}
