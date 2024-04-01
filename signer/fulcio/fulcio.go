// Copyright 2023 The Witness Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fulcio

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/log"
	"github.com/in-toto/go-witness/registry"
	"github.com/in-toto/go-witness/signer"
	"github.com/mattn/go-isatty"
	fulciopb "github.com/sigstore/fulcio/pkg/generated/protobuf"
	"github.com/sigstore/sigstore/pkg/cryptoutils"
	"github.com/sigstore/sigstore/pkg/oauthflow"
	"github.com/sigstore/sigstore/pkg/signature"
	sigo "github.com/sigstore/sigstore/pkg/signature/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/square/go-jose.v2/jwt"
)

func init() {
	signer.Register("fulcio", func() signer.SignerProvider { return New() },
		registry.StringConfigOption(
			"url",
			"Fulcio address to sign with",
			"",
			func(sp signer.SignerProvider, fulcioUrl string) (signer.SignerProvider, error) {
				fsp, ok := sp.(FulcioSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a fulcio signer provider")
				}

				WithFulcioURL(fulcioUrl)(&fsp)
				return fsp, nil
			},
		),
		registry.StringConfigOption(
			"oidc-issuer",
			"OIDC issuer to use for authentication",
			"",
			func(sp signer.SignerProvider, oidcIssuer string) (signer.SignerProvider, error) {
				fsp, ok := sp.(FulcioSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a fulcio signer provider")
				}

				WithOidcIssuer(oidcIssuer)(&fsp)
				return fsp, nil
			},
		),
		registry.StringConfigOption(
			"oidc-client-id",
			"OIDC client ID to use for authentication",
			"",
			func(sp signer.SignerProvider, oidcClientID string) (signer.SignerProvider, error) {
				fsp, ok := sp.(FulcioSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a fulcio signer provider")
				}

				WithOidcClientID(oidcClientID)(&fsp)
				return fsp, nil
			},
		),
		registry.StringConfigOption(
			"token",
			"Raw token string to use for authentication to fulcio (cannot be used in conjunction with --fulcio-token-path)",
			"",
			func(sp signer.SignerProvider, token string) (signer.SignerProvider, error) {
				fsp, ok := sp.(FulcioSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a fulcio signer provider")
				}

				WithToken(token)(&fsp)
				return fsp, nil
			},
		),
		registry.StringConfigOption(
			"oidc-redirect-url",
			"OIDC redirect URL (Optional). The default oidc-redirect-url is 'http://localhost:0/auth/callback'.",
			"",
			func(sp signer.SignerProvider, oidcRedirectUrl string) (signer.SignerProvider, error) {
				fsp, ok := sp.(FulcioSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a fulcio signer provider")
				}

				WithOidcRedirectUrl(oidcRedirectUrl)(&fsp)
				return fsp, nil
			},
		),
		registry.StringConfigOption(
			"token-path",
			"Path to the file containing a raw token to use for authentication to fulcio (cannot be used in conjunction with --fulcio-token)",
			"",
			func(sp signer.SignerProvider, tokenPath string) (signer.SignerProvider, error) {
				fsp, ok := sp.(FulcioSignerProvider)
				if !ok {
					return sp, fmt.Errorf("provided signer provider is not a fulcio signer provider")
				}

				WithTokenPath(tokenPath)(&fsp)
				return fsp, nil
			},
		),
	)
}

type FulcioSignerProvider struct {
	FulcioURL       string
	OidcIssuer      string
	OidcClientID    string
	Token           string
	TokenPath       string
	OidcRedirectUrl string
}

type Option func(*FulcioSignerProvider)

func WithFulcioURL(url string) Option {
	return func(fsp *FulcioSignerProvider) {
		fsp.FulcioURL = url
	}
}

func WithOidcIssuer(oidcIssuer string) Option {
	return func(fsp *FulcioSignerProvider) {
		fsp.OidcIssuer = oidcIssuer
	}
}

func WithOidcClientID(oidcClientID string) Option {
	return func(fsp *FulcioSignerProvider) {
		fsp.OidcClientID = oidcClientID
	}
}

func WithToken(tokenOption string) Option {
	return func(fsp *FulcioSignerProvider) {
		fsp.Token = tokenOption
	}
}

func WithOidcRedirectUrl(oidcRedirectUrl string) Option {
	return func(fsp *FulcioSignerProvider) {
		fsp.OidcRedirectUrl = oidcRedirectUrl
	}
}

func WithTokenPath(tokenPathOption string) Option {
	return func(fsp *FulcioSignerProvider) {
		fsp.TokenPath = tokenPathOption
	}
}

func New(opts ...Option) FulcioSignerProvider {
	fsp := FulcioSignerProvider{}
	for _, opt := range opts {
		opt(&fsp)
	}

	return fsp
}

func (fsp FulcioSignerProvider) Signer(ctx context.Context) (cryptoutil.Signer, error) {
	// Parse the Fulcio URL to extract its components
	u, err := url.Parse(fsp.FulcioURL)
	if err != nil {
		return nil, err
	}

	// Get the scheme, default to HTTPS if not present
	scheme := u.Scheme
	if scheme == "" {
		scheme = "https"
	}

	// Get the port, default to 443
	port := 443
	if u.Port() != "" {
		p, err := strconv.Atoi(u.Port())
		if err != nil {
			return nil, fmt.Errorf("invalid port in Fulcio URL: %s", u.Port())
		}
		port = p
	}

	// Get the host, return an error if not present
	if u.Host == "" {
		return nil, errors.New("fulcio URL must include a host")
	}

	// Make insecure true only if the scheme is HTTP
	insecure := scheme == "http"

	fClient, err := newClient(ctx, scheme+"://"+u.Host, port, insecure)
	if err != nil {
		return nil, err
	}

	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, err
	}

	var raw string

	switch {
	case fsp.Token == "" && fsp.TokenPath == "" && os.Getenv("GITHUB_ACTIONS") == "true":
		tokenURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
		if tokenURL == "" {
			return nil, errors.New("ACTIONS_ID_TOKEN_REQUEST_URL is not set")
		}

		requestToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")
		if requestToken == "" {
			return nil, errors.New("ACTIONS_ID_TOKEN_REQUEST_TOKEN is not set")
		}

		raw, err = fetchToken(tokenURL, requestToken, "sigstore")
		if err != nil {
			return nil, err
		}
	// we want to fail if both flags used (they're mutually exclusive)
	case fsp.TokenPath != "" && fsp.Token != "":
		return nil, errors.New("only one of --fulcio-token-path or --fulcio-raw-token can be used")
	case fsp.Token != "" && fsp.TokenPath == "":
		raw = fsp.Token
	case fsp.TokenPath != "" && fsp.Token == "":
		f, err := os.ReadFile(fsp.TokenPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read fulcio token from filepath %s: %w", fsp.TokenPath, err)
		}

		raw = string(f)
	case fsp.Token == "" && isatty.IsTerminal(os.Stdin.Fd()):
		tok, err := oauthflow.OIDConnect(fsp.OidcIssuer, fsp.OidcClientID, "", fsp.OidcRedirectUrl, oauthflow.DefaultIDTokenGetter)
		if err != nil {
			return nil, err
		}
		raw = tok.RawString
	default:
		return nil, errors.New("no token provided")
	}

	certResp, err := getCert(ctx, key, fClient, raw)
	if err != nil {
		return nil, err
	}

	var chain *fulciopb.CertificateChain

	switch cert := certResp.Certificate.(type) {
	case *fulciopb.SigningCertificate_SignedCertificateDetachedSct:
		chain = cert.SignedCertificateDetachedSct.GetChain()
	case *fulciopb.SigningCertificate_SignedCertificateEmbeddedSct:
		chain = cert.SignedCertificateEmbeddedSct.GetChain()
	}

	certs := chain.Certificates

	var rootCACert *x509.Certificate
	var intermediateCerts []*x509.Certificate
	var leafCert *x509.Certificate

	for _, certPEM := range certs {
		certDER, _ := pem.Decode([]byte(certPEM))
		if certDER == nil {
			return nil, errors.New("failed to decode PEM block containing the certificate")
		}

		cert, err := x509.ParseCertificate(certDER.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate: %w", err)
		}

		if cert.IsCA {
			if cert.Subject.CommonName == cert.Issuer.CommonName {
				rootCACert = cert
			} else {
				intermediateCerts = append(intermediateCerts, cert)
			}
		} else {
			leafCert = cert
		}
	}

	ss := cryptoutil.NewECDSASigner(key, crypto.SHA256)
	if ss == nil {
		return nil, errors.New("failed to create RSA signer")
	}

	witnessSigner, err := cryptoutil.NewX509Signer(ss, leafCert, intermediateCerts, []*x509.Certificate{rootCACert})
	if err != nil {
		return nil, err
	}

	return witnessSigner, nil
}

func getCert(ctx context.Context, key *ecdsa.PrivateKey, fc fulciopb.CAClient, token string) (*fulciopb.SigningCertificate, error) {
	t, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt token for fulcio: %w", err)
	}

	var claims struct {
		jwt.Claims
		Email   string `json:"email"`
		Subject string `json:"sub"`
	}

	if err := t.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return nil, err
	}

	var tok *oauthflow.OIDCIDToken

	if claims.Email != "" {
		tok = &oauthflow.OIDCIDToken{
			RawString: token,
			Subject:   claims.Email,
		}
	} else if claims.Subject != "" {
		tok = &oauthflow.OIDCIDToken{
			RawString: token,
			Subject:   claims.Subject,
		}
	}

	if tok == nil || tok.Subject == "" {
		return nil, errors.New("no email or subject claim found in token")
	}

	msg := strings.NewReader(tok.Subject)

	signer, err := signature.LoadSigner(key, crypto.SHA256)
	if err != nil {
		return nil, err
	}

	proof, err := signer.SignMessage(msg, sigo.WithCryptoSignerOpts(crypto.SHA256))
	if err != nil {
		return nil, err
	}

	pubBytesPEM, err := cryptoutils.MarshalPublicKeyToPEM(&key.PublicKey)
	if err != nil {
		return nil, err
	}

	cscr := &fulciopb.CreateSigningCertificateRequest{
		Credentials: &fulciopb.Credentials{
			Credentials: &fulciopb.Credentials_OidcIdentityToken{
				OidcIdentityToken: token,
			},
		},
		Key: &fulciopb.CreateSigningCertificateRequest_PublicKeyRequest{
			PublicKeyRequest: &fulciopb.PublicKeyRequest{
				PublicKey: &fulciopb.PublicKey{
					Content: string(pubBytesPEM),
				},
				ProofOfPossession: proof,
			},
		},
	}

	sc, err := fc.CreateSigningCertificate(ctx, cscr)
	if err != nil {
		log.Info("Failed creating signing certificate")
		return nil, err
	}

	return sc, nil
}

func newClient(ctx context.Context, fulcioURL string, fulcioPort int, isInsecure bool) (fulciopb.CAClient, error) {
	if isInsecure {
		log.Infof("Fulcio client is running in insecure mode")
	}

	// Parse the Fulcio URL
	u, err := url.Parse(fulcioURL)
	if err != nil {
		return nil, err
	}

	// Verify that the URL is valid based on the isInsecure flag
	if (u.Scheme != "https" && !isInsecure) || u.Host == "" {
		return nil, fmt.Errorf("invalid Fulcio URL: %s", fulcioURL)
	}

	// Set up the TLS configuration
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if isInsecure {
		tlsConfig.InsecureSkipVerify = true
	}

	creds := credentials.NewTLS(tlsConfig)

	// Set up the gRPC dial options
	dialOpts := []grpc.DialOption{
		grpc.WithAuthority(u.Hostname()),
	}

	if isInsecure {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	}

	// Dial the gRPC server
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(dialCtx, net.JoinHostPort(u.Hostname(), strconv.Itoa(fulcioPort)), dialOpts...)
	if err != nil {
		return nil, err
	}

	// Create the Fulcio client
	return fulciopb.NewCAClient(conn), nil
}
