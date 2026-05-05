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

func Create(payload Payload, writer http.ResponseWriter) {
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
	/* Major screw up */
	if err != nil {
		panic(err)
	}
	cookies.Set(writer, tokenStr, payload.RememberMe)
}

func CreateFresh(username string, rememberMe bool, writer http.ResponseWriter) {
	fromDate, toDate := utils.GetMonthBounds()
	payload := Payload{
		UserID:     username,
		RememberMe: rememberMe,
		PageSize:   config.DefaultPageSize,
		Page:       1,
		SearchBy:   "",
		SortBy:     1,
		SortAsc:    true,
		FromDate:   fromDate.Format("2006-01-02"),
		ToDate:     toDate.Format("2006-01-02"),
	}
	Create(payload, writer)
}

func Update(payload *Payload, key string, value any, writer http.ResponseWriter) {
	reflectValue := reflect.ValueOf(payload).Elem()
	field := reflectValue.FieldByName(key)
	if !field.IsValid() || !field.CanSet() {
		panic(fmt.Errorf("Invalid payload key: %s", key))
	}
	switch key {
	case "UserID", "SearchBy", "FromDate", "ToDate":
		valueStr, ok := value.(string)
		if !ok {
			panic("Invalid payload type (must be string)")
		}
		field.SetString(valueStr)
	case "Page", "PageSize":
		valueInt, ok := value.(int)
		if !ok {
			panic("Invalid payload type (must be int)")
		}
		field.SetInt(int64(valueInt))
	case "SortBy":
		valueInt, ok := value.(utils.SortableColumn)
		if !ok {
			panic("Invalid payload type (must be SortableColumn)")
		}
		field.SetInt(int64(utils.SortableColumn(valueInt)))
	case "RememberMe", "SortAsc":
		valueBool, ok := value.(bool)
		if !ok {
			panic("Invalid payload type (must be bool)")
		}
		field.SetBool(valueBool)
	default:
		panic(fmt.Errorf("Invalid payload key: %s", key))
	}
	Create(*payload, writer)
}

func Get(req *http.Request) *Payload {
	cookie := cookies.Get(req)
	if cookie == "" {
		panic("Empty cookie")
	}
	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			panic(fmt.Errorf("Unexpected signing method: %v", token.Header["alg"]))
		}
		return config.JWTSecret, nil
	})
	if err != nil {
		panic(fmt.Errorf("Unable to parse token: %w", err))
	}
	if !parsed.Valid {
		panic("Invalid token")
	}
	return &claims.Payload
}
