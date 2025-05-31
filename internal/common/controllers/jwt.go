package controllers

import (
	"accounts/internal/common/logger"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTController gestiona generación y validación de JWTs usando RSA (RS256).
// PrivateKey y PublicKey deben contener las claves en formato PEM.
type JWTController struct {
	PrivateKey string
	PublicKey  string
}

func normalizePEM(pemStr string) string {
	// Si ya tiene saltos de línea, asumimos que está normalizado
	if strings.Contains(pemStr, "\n") {
		return pemStr
	}

	// Intentamos detectar si es una línea larga con espacios
	parts := strings.Split(pemStr, " ")
	if len(parts) < 3 {
		// Probablemente ya esté bien o está corrupto, lo regresamos igual
		return pemStr
	}

	header := parts[0]
	footer := parts[len(parts)-1]
	body := strings.Join(parts[1:len(parts)-1], "")

	return fmt.Sprintf("%s\n%s\n%s", header, body, footer)
}

// parseRSAPrivateKey decodifica un PEM PKCS#1 o PKCS#8 en *rsa.PrivateKey.
func parseRSAPrivateKey(ctx context.Context, pemStr string) (*rsa.PrivateKey, error) {
	entry := logger.FromContext(ctx)
	entry.Info("Parsing RSA private key")

	pemStrNormalized := normalizePEM(pemStr)
	//pemStrNormalized := pemStr

	block, _ := pem.Decode([]byte(pemStrNormalized))
	if block == nil {
		entry.Error("Failed to decode PEM private key")
		return nil, errors.New("no se pudo decodificar PEM de clave privada")
	}
	// PKCS#1
	if priv, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		entry.Info("RSA private key parsed successfully")
		return priv, nil
	}
	// PKCS#8
	keyIfc, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		entry.Error("Failed to parse PKCS#8 private key", "error", err)
		return nil, fmt.Errorf("falló parseo PKCS#8: %w", err)
	}
	priv, ok := keyIfc.(*rsa.PrivateKey)
	if !ok {
		entry.Error("PKCS#8 private key is not RSA")
		return nil, errors.New("clave PKCS#8 no es RSA")
	}
	entry.Info("RSA private key parsed successfully")
	return priv, nil
}

// parseRSAPublicKey decodifica un PEM en *rsa.PublicKey.
func parseRSAPublicKey(ctx context.Context, pemStr string) (*rsa.PublicKey, error) {
	entry := logger.FromContext(ctx)
	entry.Info("Parsing RSA public key")

	pemStrNormalized := normalizePEM(pemStr)
	//pemStrNormalized := pemStr

	block, _ := pem.Decode([]byte(pemStrNormalized))
	if block == nil {
		entry.Error("Failed to decode PEM public key")
		return nil, errors.New("no se pudo decodificar PEM de clave pública")
	}
	pubIfc, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		entry.Error("Failed to parse public key", "error", err)
		return nil, fmt.Errorf("falló parseo clave pública: %w", err)
	}
	pub, ok := pubIfc.(*rsa.PublicKey)
	if !ok {
		entry.Error("Public key is not RSA")
		return nil, errors.New("clave pública no es RSA")
	}
	entry.Info("RSA public key parsed successfully")
	return pub, nil
}

// GenerateToken crea un JWT con claims personalizados y lo firma con RS256.
func (j *JWTController) GenerateToken(ctx context.Context, data map[string]interface{}, expireSeconds int) (string, error) {

	entry := logger.FromContext(ctx)
	entry.Info("Generating JWT token")

	privKey, err := parseRSAPrivateKey(ctx, j.PrivateKey)
	if err != nil {
		entry.Error("Error parsing private key", "error", err)
		return "", err
	}
	entry.Info("Private key parsed successfully")
	claims := jwt.MapClaims{}
	for k, v := range data {
		claims[k] = v
	}
	if expireSeconds > 0 {
		entry.Info("Adding expiration time to claims")
		claims["exp"] = time.Now().Add(time.Duration(expireSeconds) * time.Second).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	entry.Info("Token created successfully")
	return token.SignedString(privKey)
}

// ValidateToken valida un JWT firmado con RS256 y retorna sus claims.
// Si el token ha expirado, devuelve un error con mensaje personalizado.
func (j *JWTController) ValidateToken(ctx context.Context, tokenString string) (map[string]interface{}, error) {
	entry := logger.FromContext(ctx)
	entry.Info("Validating JWT token")

	pubKey, err := parseRSAPublicKey(ctx, j.PublicKey)
	if err != nil {
		entry.Error("Error parsing public key", "error", err)
		return nil, err
	}

	entry.Info("Public key parsed successfully")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			entry.Error("Unexpected algorithm", "algorithm", token.Method.Alg())
			return nil, fmt.Errorf("algoritmo inesperado: %s", token.Method.Alg())
		}
		return pubKey, nil
	})
	entry.Info("Token parsed successfully")
	// Si hay un error al parsear, comprobamos si se debe a expiración
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) && ve.Errors&jwt.ValidationErrorExpired != 0 {
			entry.Error("Token expired", "error", ve.Inner)
			return nil, fmt.Errorf("%v", ve.Inner)
		}
		entry.Error("Error parsing token", "error", err)
		return nil, err
	}

	// Si el token es válido, retornamos las claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		entry.Info("Token validated successfully")
		return claims, nil
	}
	entry.Error("Token invalid")
	return nil, errors.New("token inválido")
}
