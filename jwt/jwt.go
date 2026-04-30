package jwt

import (
	"fmt"
	"net/http"
	"reflect"
	"time"
	"todoer/config"
	"todoer/cookies"
	"todoer/utils"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	UserID     string               `json:"user_id"`
	RememberMe bool                 `json:"remember_me"`
	SearchBy   string               `json:"search_by"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	SortBy     utils.SortableColumn `json:"sort_by"`
	SortAsc    bool                 `json:"sort_asc"`
	FromDate   string               `json:"from_date"`
	ToDate     string               `json:"to_date"`
}

type Claims struct {
	Payload
	jwt.RegisteredClaims
}

func HealthCheck(payload *Payload, page int, writer http.ResponseWriter) {
	if payload.Page != page {
		payload.Page = page
		token := Create(*payload)
		cookies.Set(writer, token, payload.RememberMe)
	}
}

func Create(payload Payload) string {
	expirationTime := time.Now()
	if payload.RememberMe {
		expirationTime = expirationTime.Add(
			time.Duration(config.CookieLifetime) * time.Second,
		)
	} else {
		expirationTime = expirationTime.Add(
			time.Duration(config.CookieShortLifetime) * time.Second,
		)
	}
	claims := Claims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(config.JWTSecret)
	/* Something went terribly wrong */
	if err != nil {
		panic(err)
	}
	return tokenStr
}

func Update(payload *Payload, key string, value any, writer http.ResponseWriter) error {
	reflectValue := reflect.ValueOf(payload).Elem()
	field := reflectValue.FieldByName(key)
	if !field.IsValid() || !field.CanSet() {
		return fmt.Errorf("Invalid payload key: %s", key)
	}
	switch key {
	case "UserID", "SearchBy", "FromDate", "ToDate":
		valueStr, ok := value.(string)
		if !ok {
			return fmt.Errorf("Invalid payload type (must be string)")
		}
		field.SetString(valueStr)
	case "Page", "PageSize":
		valueInt, ok := value.(int)
		if !ok {
			return fmt.Errorf("Invalid payload type (must be int)")
		}
		field.SetInt(int64(valueInt))
	case "SortBy":
		valueInt, ok := value.(int)
		if !ok {
			return fmt.Errorf("Invalid payload type (must be int)")
		}
		field.SetInt(int64(utils.SortableColumn(valueInt)))
	case "RememberMe", "SortAsc":
		valueBool, ok := value.(bool)
		if !ok {
			return fmt.Errorf("Invalid payload type (must be bool)")
		}
		field.SetBool(valueBool)
	default:
		return fmt.Errorf("Invalid payload key: %s", key)
	}
	token := Create(*payload)
	cookies.Set(writer, token, payload.RememberMe)
	return nil
}

func GetPayload(tokenStr string) (*Payload, error) {
	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return config.JWTSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to parse token: %w", err)
	}
	if !parsed.Valid {
		return nil, fmt.Errorf("Invalid token")
	}
	return &claims.Payload, nil
}

func Get(req *http.Request) (*Payload, error) {
	cookie := cookies.Get(req)
	if cookie == "" {
		return nil, fmt.Errorf("Empty cookie")
	}
	payload, err := GetPayload(cookie)
	if err != nil {
		return nil, fmt.Errorf("Unable to read token: %w", err)
	}
	return payload, nil
}

func Set(writer http.ResponseWriter, payload *Payload, rememberMe bool) {
	token := Create(*payload)
	cookies.Set(writer, token, rememberMe)
}
