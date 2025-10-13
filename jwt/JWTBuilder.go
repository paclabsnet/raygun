/*
Copyright © 2025 PACLabs
*/
package jwt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"raygun/config"
	"raygun/log"
	"raygun/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTBuilder struct {
}

func NewJWTBuilder() JWTBuilder {
	builder := JWTBuilder{}
	return builder
}

/*
 * First, we create all the claims. Some will be from config data, others (exp) will be derived
 * Then, we
 */
func (builder JWTBuilder) Generate(suiteConfig types.TestSuite, jwtConfig types.TestJwt) (string, error) {

	log.Debug("JWTBuilder: SuiteConfig: %v, testconfig: %v", suiteConfig.Jwt, jwtConfig)

	// Create claims map
	claims, err := createClaims(suiteConfig.Jwt, jwtConfig)

	if err != nil {
		return "", err
	}

	// Get signing method and key
	signingMethod := jwt.GetSigningMethod(choose(suiteConfig.Jwt.Algorithm, jwtConfig.Algorithm))
	if signingMethod == nil {
		return "", fmt.Errorf("unsupported signing algorithm: %s/%s", suiteConfig.Jwt.Algorithm, jwtConfig.Algorithm)
	}

	// Create token
	token := jwt.NewWithClaims(signingMethod, claims)

	// Sign token with appropriate key
	signingKey, err := getSigningKey(suiteConfig.Jwt, jwtConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get signing key: %w", err)
	}

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil

}

/*
 *  Create the JWT claims - most are pulled from the config data, but
 *  some may be time-sensitive, and created via calculation
 */
func createClaims(suiteConfig types.TestJwt, jwtConfig types.TestJwt) (jwt.MapClaims, error) {

	claims := jwt.MapClaims{}

	// Add standard claims
	if jwtConfig.Claims.Issuer != "" {
		claims["iss"] = config.Resolver.ExpandProperties(choose(suiteConfig.Claims.Issuer, jwtConfig.Claims.Issuer))
	}
	if jwtConfig.Claims.Subject != "" {
		claims["sub"] = config.Resolver.ExpandProperties(choose(suiteConfig.Claims.Subject, jwtConfig.Claims.Subject))
	}
	if jwtConfig.Claims.Audience != nil {
		claims["aud"] = jwtConfig.Claims.Audience
	}
	if jwtConfig.Claims.JWTID != "" {
		claims["jti"] = config.Resolver.ExpandProperties(choose(suiteConfig.Claims.JWTID, jwtConfig.Claims.JWTID))
	}

	// Handle time-based claims
	now := time.Now()
	if jwtConfig.Claims.IncludeIat {
		claims["iat"] = now.Unix()
	}

	lifetime := choose(suiteConfig.Claims.Lifetime, jwtConfig.Claims.Lifetime)

	if lifetime == "" {
		lifetime = "3600s"
	}

	duration, err := time.ParseDuration(lifetime)
	if err != nil {
		return nil, fmt.Errorf("invalid lifetime, should be a golang duration, like  3600s or 2h or 180000ms")
	}

	claims["exp"] = now.Add(duration).Unix()

	// process the suite default settings, and then allow for
	// test-specific overrides
	for key, value := range suiteConfig.Claims.Custom {
		if strVal, ok := value.(string); ok {
			claims[key] = config.Resolver.ExpandProperties(strVal)
		} else {
			claims[key] = value
		}
	}

	for key, value := range jwtConfig.Claims.Custom {
		if strVal, ok := value.(string); ok {
			claims[key] = config.Resolver.ExpandProperties(strVal)
		} else {
			claims[key] = value
		}
	}

	return claims, nil
}

// getSigningKey returns the appropriate signing key based on algorithm
func getSigningKey(suiteConfig types.TestJwt, jwtConfig types.TestJwt) (interface{}, error) {

	algorithm := config.Resolver.ExpandProperties(choose(suiteConfig.Algorithm, jwtConfig.Algorithm))

	// HMAC algorithms use secret key
	if algorithm == "HS256" || algorithm == "HS384" || algorithm == "HS512" {
		return getSecretKey(choose(suiteConfig.Secret, jwtConfig.Secret))
	}

	// RSA, ECDSA, and EdDSA algorithms use private keys.
	if jwtConfig.PrivateKey == "" {
		return nil, fmt.Errorf("private_key_path is required for %s algorithm", algorithm)
	}

	// for security, we allow the private key to be pulled from an environment
	// variable
	keyData := config.Resolver.ExpandProperties(choose(suiteConfig.PrivateKey, jwtConfig.PrivateKey))

	return parsePrivateKey([]byte(keyData), algorithm)
}

/**
 *  the secret key is a shared secret. It may be stored in an environment
 *  variable, so we can pull it from there
 */
func getSecretKey(key string) (interface{}, error) {

	if key == "" {
		return nil, fmt.Errorf("secret is required for HMAC algorithms")
	}

	expanded_key := config.Resolver.ExpandProperties(key)

	log.Debug("ExpandedKey: [%s]", expanded_key)
	return []byte(expanded_key), nil

}

// parsePrivateKey parses a PEM-encoded private key
func parsePrivateKey(keyData []byte, algorithm string) (interface{}, error) {
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Try parsing based on algorithm
	switch {
	case algorithm == "RS256" || algorithm == "RS384" || algorithm == "RS512" ||
		algorithm == "PS256" || algorithm == "PS384" || algorithm == "PS512":
		// RSA keys
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			// Try PKCS8 format
			keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
			}
			rsaKey, ok := keyInterface.(*rsa.PrivateKey)
			if !ok {
				return nil, fmt.Errorf("not an RSA private key")
			}
			return rsaKey, nil
		}
		return key, nil

	case algorithm == "ES256" || algorithm == "ES384" || algorithm == "ES512":
		// ECDSA keys
		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			// Try PKCS8 format
			keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse ECDSA private key: %w", err)
			}
			ecKey, ok := keyInterface.(*ecdsa.PrivateKey)
			if !ok {
				return nil, fmt.Errorf("not an ECDSA private key")
			}
			return ecKey, nil
		}
		return key, nil

	case algorithm == "EdDSA":
		// EdDSA (Ed25519) keys
		keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse EdDSA private key: %w", err)
		}
		edKey, ok := keyInterface.(ed25519.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not an Ed25519 private key")
		}
		return edKey, nil

	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

func choose(a string, b string) string {

	log.Debug("Choose: a:[%s] b:[%s]", a, b)

	if b != "" {
		return b
	}
	return a
}
