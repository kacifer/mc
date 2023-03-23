package mjwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)

const DefaultLease = 30 * 24 * time.Hour

const IDKey = "id"

const NameKey = "name"

type Engine interface {
	SignForID(id uint) (token *jwt.Token)
	SignedStringForID(id uint) (tokenString string, err error)
	SignForName(name string) (token *jwt.Token)
	SignedStringForName(name string) (tokenString string, err error)
	Parse(tokenString string) (*jwt.Token, error)
	Validate(token *jwt.Token) (jwt.MapClaims, error)
	ValidateSignedString(tokenString string) (*jwt.Token, jwt.MapClaims, error)
	ValidateHeader(authHeader string) (*jwt.Token, jwt.MapClaims, error)
	ExtractID(token *jwt.Token) (id uint, err error)
	ExtractIDFromSignedString(tokenString string) (id uint, err error)
	ExtractIDFromHeader(authHeader string) (id uint, err error)
	ExtractName(token *jwt.Token) (name string, err error)
	ExtractNameFromSignedString(tokenString string) (name string, err error)
	ExtractNameFromHeader(authHeader string) (name string, err error)
}

type EngineImpl struct {
	Secret []byte

	Lease time.Duration

	NowFunc func() time.Time
}

func (e *EngineImpl) SignMapClaims(claims jwt.MapClaims) (token *jwt.Token) {
	claims["iat"] = e.NowFunc()
	claims["exp"] = e.NowFunc().Add(e.Lease)
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// SignForID and get the complete encoded token as a string using the secret
	return token
}

// SignForID and get the complete encoded token as a string using the secret
func (e *EngineImpl) SignForID(id uint) (token *jwt.Token) {
	return e.SignMapClaims(jwt.MapClaims{
		IDKey: id,
	})
}

func (e *EngineImpl) SignedStringForID(id uint) (tokenString string, err error) {
	return e.SignForID(id).SignedString(e.Secret)
}

// SignForName and get the complete encoded token as a string using the secret
func (e *EngineImpl) SignForName(name string) (token *jwt.Token) {
	return e.SignMapClaims(jwt.MapClaims{
		NameKey: name,
	})
}

func (e *EngineImpl) SignedStringForName(name string) (tokenString string, err error) {
	return e.SignForName(name).SignedString(e.Secret)
}

// Parse do not Validate the token payload
func (e *EngineImpl) Parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return e.Secret, nil
	})
}

func (e *EngineImpl) Validate(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.Errorf("invalid claims")
	}
	if !token.Valid {
		return nil, errors.Errorf("invalid token")
	}
	return claims, nil
}

func (e *EngineImpl) ValidateSignedString(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := e.Parse(tokenString)
	if err != nil {
		return nil, nil, errors.Wrap(err, "parse error")
	}
	claims, err := e.Validate(token)
	if err != nil {
		return nil, nil, errors.Wrap(err, "validate error")
	}
	return token, claims, nil
}

func (e *EngineImpl) ValidateHeader(authHeader string) (*jwt.Token, jwt.MapClaims, error) {
	if authHeader == "" {
		return nil, nil, errors.New("auth header missing")
	}
	if len(authHeader) < 7 {
		return nil, nil, errors.Errorf("malformed auth header")
	}
	return e.ValidateSignedString(authHeader[7:])
}

func (e *EngineImpl) ExtractIDFromClaims(claims jwt.MapClaims) (id uint, err error) {
	switch claims[IDKey].(type) {
	case float64:
		id = uint(claims[IDKey].(float64))
	case uint:
		id = claims[IDKey].(uint)
	default:
		var ok bool
		id, ok = claims[IDKey].(uint)
		if !ok {
			return 0, errors.Errorf("invalid uid found")
		}
	}
	return id, nil
}

func (e *EngineImpl) ExtractID(token *jwt.Token) (id uint, err error) {
	claims, err := e.Validate(token)
	if err != nil {
		return 0, err
	}
	return e.ExtractIDFromClaims(claims)
}

func (e *EngineImpl) ExtractIDFromSignedString(tokenString string) (id uint, err error) {
	_, claims, err := e.ValidateSignedString(tokenString)
	if err != nil {
		return 0, err
	}
	return e.ExtractIDFromClaims(claims)
}

func (e *EngineImpl) ExtractIDFromHeader(authHeader string) (id uint, err error) {
	_, claims, err := e.ValidateHeader(authHeader)
	if err != nil {
		return 0, err
	}
	return e.ExtractIDFromClaims(claims)
}

func (e *EngineImpl) ExtractNameFromClaims(claims jwt.MapClaims) (name string, err error) {
	switch claims[NameKey].(type) {
	case string:
		name = claims[NameKey].(string)
	default:
		var ok bool
		name, ok = claims[NameKey].(string)
		if !ok {
			return "", errors.Errorf("invalid name found")
		}
	}
	return name, nil
}

func (e *EngineImpl) ExtractName(token *jwt.Token) (name string, err error) {
	claims, err := e.Validate(token)
	if err != nil {
		return "", err
	}
	return e.ExtractNameFromClaims(claims)
}

func (e *EngineImpl) ExtractNameFromSignedString(tokenString string) (name string, err error) {
	_, claims, err := e.ValidateSignedString(tokenString)
	if err != nil {
		return "", err
	}
	return e.ExtractNameFromClaims(claims)
}

func (e *EngineImpl) ExtractNameFromHeader(authHeader string) (name string, err error) {
	_, claims, err := e.ValidateHeader(authHeader)
	if err != nil {
		return "", err
	}
	return e.ExtractNameFromClaims(claims)
}

func NewImpl(secret []byte, lease time.Duration) *EngineImpl {
	return &EngineImpl{
		Secret:  secret,
		Lease:   lease,
		NowFunc: time.Now,
	}
}

func New(secret []byte, lease time.Duration) Engine {
	return NewImpl(secret, lease)
}

func NewDefault(secret []byte) Engine {
	return New(secret, DefaultLease)
}
