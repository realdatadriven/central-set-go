// internal/auth/oidc.go
import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"sync"
)

type ProviderJWKS struct {
	Set     jwk.Set
	Mutex   sync.RWMutex
	KeyFunc jwt.Keyfunc
}

var providerJWKS = make(map[string]*ProviderJWKS) // providerKey → JWKS

func LoadJWKSForProvider(providerKey, jwksURL string) error {
	set, err := jwk.Fetch(context.Background(), jwksURL)
	if err != nil {
		return err
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing kid")
		}

		key, found := set.LookupKeyID(kid)
		if !found {
			// Optional: refresh JWKS on miss (cache invalidation)
			return nil, fmt.Errorf("unknown kid")
		}

		var raw any
		if err := key.Raw(&raw); err != nil {
			return nil, err
		}
		return raw, nil
	}

	providerJWKS[providerKey] = &ProviderJWKS{
		Set:     set,
		KeyFunc: keyFunc,
	}
	return nil
}

// ValidateIDToken example
func ValidateIDToken(providerKey, idToken string) (*jwt.MapClaims, error) {
	p, ok := providerJWKS[providerKey]
	if !ok {
		return nil, fmt.Errorf("no JWKS for provider %s", providerKey)
	}

	token, err := jwt.ParseWithClaims(idToken, jwt.MapClaims{}, p.KeyFunc,
		jwt.WithIssuer(getExpectedIssuer(providerKey)),
		jwt.WithAudience(getExpectedAudience(providerKey)),
		jwt.WithLeeway(5*time.Second), // small clock skew tolerance
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid id_token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return &claims, nil
}