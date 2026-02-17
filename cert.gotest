package main

/*
# Generate CA
openssl genrsa -out ca-key.pem 4096
openssl req -new -x509 -days 365 -key ssl/ca-key.pem -out ssl/ca-cert.pem -subj "/CN=Test CA"

# Generate server certificate
openssl genrsa -out ssl/server-key.pem 4096
openssl req -new -key ssl/server-key.pem -out ssl/server.csr -subj "/CN=localhost"
openssl x509 -req -days 365 -in ssl/server.csr -CA ssl/ca-cert.pem -CAkey ssl/ca-key.pem -CAcreateserial -out ssl/server-cert.pem
*/

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
)

func main() {
	// ───────────────────────────────────────────────
	//  1. Generate CA key + self-signed CA certificate
	// ───────────────────────────────────────────────
	caKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal("failed to generate CA key:", err)
	}

	caTemplate := &x509.Certificate{
		SerialNumber:          randomSerial(),
		Subject:               pkix.Name{CommonName: "Test CA"},
		NotBefore:             time.Now().Add(-24 * time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}

	caCertDER, err := x509.CreateCertificate(
		rand.Reader,
		caTemplate,
		caTemplate, // self-signed
		&caKey.PublicKey,
		caKey,
	)
	if err != nil {
		log.Fatal("failed to create CA cert:", err)
	}

	// Save CA key
	saveKey("ssl/ca-key.pem", caKey)
	// Save CA cert
	saveCert("ssl/ca-cert.pem", caCertDER)

	fmt.Println("Created CA certificate and key")

	// ───────────────────────────────────────────────
	//  2. Generate server key + certificate signed by CA
	// ───────────────────────────────────────────────
	serverKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal("failed to generate server key:", err)
	}

	serverTemplate := &x509.Certificate{
		SerialNumber: randomSerial(),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now().Add(-24 * time.Hour),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,

		// You can add these if needed:
		// DNSNames:       []string{"localhost", "127.0.0.1"},
		// IPAddresses:    []net.IP{net.ParseIP("127.0.0.1")},
	}

	serverCertDER, err := x509.CreateCertificate(
		rand.Reader,
		serverTemplate,
		caTemplate,           // signed by CA
		&serverKey.PublicKey, // public key of server
		caKey,                // CA private key
	)
	if err != nil {
		log.Fatal("failed to create server cert:", err)
	}

	// Save server key + cert
	saveKey("ssl/server-key.pem", serverKey)
	saveCert("ssl/server-cert.pem", serverCertDER)

	fmt.Println("Created server certificate and key")
	fmt.Println("Done. Files written to ssl/ directory")
}

// Helpers ────────────────────────────────────────────────

func randomSerial() *big.Int {
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 159))
	if err != nil {
		log.Fatal("failed to generate serial:", err)
	}
	return serial
}

func saveKey(filename string, key *rsa.PrivateKey) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
}

func saveCert(filename string, derBytes []byte) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	pem.Encode(f, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
}
