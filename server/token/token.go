package token

import (
	"fmt"
	"net/http"
	"reflect"
	"time"
	"todoer/config"
	"todoer/utils"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	UserID     string              `json:"user_id"`
	RememberMe bool                `json:"remember_me"`
	SearchBy   string              `json:"search_by"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	SortBy     utils.SortableField `json:"sort_by"`
	SortAsc    bool                `json:"sort_asc"`
	FromDate   string              `json:"from_date"`
	ToDate     string              `json:"to_date"`
}

type Claims struct {
	Payload
	Extra string `json:"extra"`
	jwt.RegisteredClaims
}

func getExpirationTime(rememberMe bool) time.Time {
	expires := time.Now()
	if rememberMe {
		expires = expires.Add(
			time.Duration(config.CookieLifetime) * time.Second,
		)
	} else {
		expires = expires.Add(
			time.Duration(config.CookieShortLifetime) * time.Second,
		)
	}
	return expires
}

func Create(payload Payload, writer http.ResponseWriter) {
	expires := getExpirationTime(payload.RememberMe)
	claims := Claims{
		Payload: payload,
		Extra: "Hello!",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(config.JWTSecret)
	/* Major screwup */
	if err != nil {
		panic(err)
	}
	setCookie(config.CookieName, tokenStr, expires, writer)
}

func CreateFresh(username string, rememberMe bool, writer http.ResponseWriter) {
	fromDate, toDate := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
	payload := Payload{
		UserID:     username,
		RememberMe: rememberMe,
		PageSize:   config.DefaultPageSize,
		Page:       1,
		SearchBy:   "",
		SortBy:     utils.Datetime,
		SortAsc:    true,
		FromDate:   fromDate.Format(utils.HTMLDateFormat),
		ToDate:     toDate.Format(utils.HTMLDateFormat),
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
		valueInt, ok := value.(utils.SortableField)
		if !ok {
			panic("Invalid payload type (must be SortableField)")
		}
		field.SetInt(int64(utils.SortableField(valueInt)))
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
	cookie := getCookie(config.CookieName, req)
	/* Should not happen */
	if cookie == "" {
		panic("Empty cookie")
	}
	claims := Claims{}
	_, err := jwt.ParseWithClaims(cookie, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			panic(fmt.Errorf("Unexpected signing method: %v", token.Header["alg"]))
		}
		return config.JWTSecret, nil
	})
	if err != nil {
		panic(fmt.Errorf("Unable to parse token: %w", err))
	}
	return &claims.Payload
}

func Clear(writer http.ResponseWriter) {
	clearCookie(config.CookieName, writer)
}
