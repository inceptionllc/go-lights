package security

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"strings"
)

// LoadTLSCertificates loads any number of certificates from the file system
// with the given path/file name. The certificate itself should be stored in a
// `.pem` file and the private key in a `.key` file.
func LoadTLSCertificates(paths ...string) ([]tls.Certificate, error) {
	certs := []tls.Certificate{}
	for _, path := range paths {

		caRaw, err := ioutil.ReadFile(path + ".pem")
		if err != nil {
			return nil, err
		}
		keyRaw, err := ioutil.ReadFile(path + ".key")
		if err != nil {
			return nil, err
		}
		key, err := x509.ParsePKCS1PrivateKey(keyRaw)
		if err != nil {
			return nil, err
		}

		certs = append(certs, tls.Certificate{
			Certificate: [][]byte{caRaw},
			PrivateKey:  key,
		})
	}
	return certs, nil
}

// LoadCertPool creates a certificate pool and loads any number of certificates
// from the file system. The certificate file names are assumed to end in `.pem`
// and if a given name does not have `.pem` it will be added automatically.
func LoadCertPool(paths ...string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	for _, path := range paths {
		if !strings.HasSuffix(path, ".pem") {
			path += ".pem"
		}
		raw, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		cert, err := x509.ParseCertificate(raw)
		if err != nil {
			return nil, err
		}
		pool.AddCert(cert)
	}
	return pool, nil
}

// LoadTLSKeyPairs loads zero or more certificate key pairs (compatible with mqtt).
func LoadTLSKeyPairs(paths ...string) ([]tls.Certificate, error) {
	certs := []tls.Certificate{}

	for _, path := range paths {
		cert, err := tls.LoadX509KeyPair(path+".pem", path+".key")
		if err != nil {
			log.Println("Load keypair error", err)
			return nil, err
		}

		// Just to print out the client certificate..
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			log.Println("parse cert error", err)
			return nil, err
		}
		certs = append(certs, cert)
	}
	return certs, nil
}

// LoadPEMCertPool creates a certificate pool and loads any number of certificates
// from the file system. The certificate names are assumed to end in `.pem`.
func LoadPEMCertPool(paths ...string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	for _, path := range paths {
		if !strings.HasSuffix(path, ".pem") {
			path += ".pem"
		}
		raw, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		pool.AppendCertsFromPEM(raw)
	}
	return pool, nil
}
