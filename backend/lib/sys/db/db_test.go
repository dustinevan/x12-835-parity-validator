package db

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// TestConnStringTLS parses each mode's connection string with pgx and checks
// the TLS settings pgx will connect with. PGSSLROOTCERT names a missing file
// throughout: if the environment reached pgx, parsing would fail.
func TestConnStringTLS(t *testing.T) {
	t.Setenv("PGSSLROOTCERT", filepath.Join(t.TempDir(), "missing.pem"))
	rootCert := writeCertificateAuthority(t)

	tests := []struct {
		mode       string
		tls        bool
		skipVerify bool
		serverName string
		roots      bool
	}{
		{TLSDisable, false, false, "", false},
		{TLSRequire, true, true, "db.example.internal", false},
		{TLSVerifyFull, true, false, "db.example.internal", true},
	}
	for _, test := range tests {
		cfg := Config{User: "u", Password: "p", Host: "db.example.internal:5432", DBName: "validator", TLSMode: test.mode, TLSRootCert: rootCert}
		parsed, err := pgconn.ParseConfig(connString(cfg))
		if err != nil {
			t.Fatalf("ParseConfig(%s) error: %v", test.mode, err)
		}
		tlsConfig := parsed.TLSConfig
		if (tlsConfig != nil) != test.tls {
			t.Fatalf("%s: TLS = %v, want %v", test.mode, tlsConfig != nil, test.tls)
		}
		if tlsConfig == nil {
			continue
		}
		if tlsConfig.InsecureSkipVerify != test.skipVerify {
			t.Errorf("%s: InsecureSkipVerify = %v, want %v", test.mode, tlsConfig.InsecureSkipVerify, test.skipVerify)
		}
		if tlsConfig.ServerName != test.serverName {
			t.Errorf("%s: ServerName = %q, want %q", test.mode, tlsConfig.ServerName, test.serverName)
		}
		if (tlsConfig.RootCAs != nil) != test.roots {
			t.Errorf("%s: root certificates set = %v, want %v", test.mode, tlsConfig.RootCAs != nil, test.roots)
		}
	}
}

// TestOpenTLSSettings checks the settings Open refuses before connecting.
func TestOpenTLSSettings(t *testing.T) {
	tests := []struct {
		mode     string
		rootCert string
		want     string
	}{
		{"", "", "mode must be disable, require, or verify-full"},
		{"verify-ca", "", "mode must be disable, require, or verify-full"},
		{TLSVerifyFull, "", "verify-full needs a root certificate file"},
		{TLSVerifyFull, filepath.Join(t.TempDir(), "missing.pem"), "unable to read CA file"},
	}
	for _, test := range tests {
		cfg := Config{User: "u", Password: "secret", Host: "localhost:5432", DBName: "validator", TLSMode: test.mode, TLSRootCert: test.rootCert}
		_, err := Open(context.Background(), cfg)
		if err == nil || !strings.Contains(err.Error(), test.want) {
			t.Errorf("Open(%q, %q) = %v, want an error containing %q", test.mode, test.rootCert, err, test.want)
		}
		if err != nil && strings.Contains(err.Error(), "secret") {
			t.Errorf("Open(%q, %q) error shows the password: %v", test.mode, test.rootCert, err)
		}
	}
}

// writeCertificateAuthority writes a self-signed CA certificate as PEM and
// returns its path.
func writeCertificateAuthority(t *testing.T) string {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "ca.pem")
	if err := os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}
