// internal/auth/jwt_service.go

package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// AccountClaims representa os dados que armazenamos no JWT
type AccountClaims struct {
	AccountID string `json:"account_id"`
	jwt.RegisteredClaims
}

// GenerateJWT cria um token válido por 7 dias
func GenerateJWT(accountID string) (string, error) {
	claims := AccountClaims{
		AccountID: accountID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateJWT valida e retorna os claims do token
func ValidateJWT(tokenStr string) (*AccountClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccountClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccountClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}

// IsAdminByToken verifica se o usuário autenticado é admin
func IsAdminByToken(account *entity.Account) bool {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminWhatsApp := os.Getenv("ADMIN_WHATSAPP")
	return account.Email == adminEmail || account.WhatsApp == adminWhatsApp
}

// ExtractAccountIDFromHeader extrai o ID da conta a partir do header Authorization
func ExtractAccountIDFromHeader(r *http.Request) (uuid.UUID, error) {
	tokenStr := ExtractTokenFromHeader(r)
	if tokenStr == "" {
		return uuid.Nil, errors.New("token não encontrado")
	}

	claims, err := ValidateJWT(tokenStr)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(claims.AccountID)
}

// ExtractTokenFromHeader pega o token do header Authorization
func ExtractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}
