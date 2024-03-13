package services

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	log "github.com/sirupsen/logrus"
	"strings"
)

const adminScope = "admin"
const tokenScopes = "scope"

type AuthService interface {
	IsAValidToken(ctx context.Context, token string) bool
	IsAnAdminToken(ctx context.Context, token string) bool
}

type keycloakAuthService struct {
	client       *gocloak.GoCloak
	clientId     string
	clientSecret string
	realm        string
}

func NewKeycloakAuthService(kcConfig *config.KeycloakConfiguration) AuthService {
	client := gocloak.NewClient(kcConfig.KeycloakUrl)
	return &keycloakAuthService{client, kcConfig.KeycloakClient, kcConfig.KeycloakClientSecret, kcConfig.KeycloakRealm}
}

func (k keycloakAuthService) IsAValidToken(ctx context.Context, token string) bool {
	rptResult, err := k.client.RetrospectToken(ctx, token, k.clientId, k.clientSecret, k.realm)
	if err != nil {
		log.Errorf("Error validating token: %v", err)
		return false
	}

	if !*rptResult.Active {
		log.Warnf("Recived a token that is not valid")
		return false
	}
	
	return true
}

func (k keycloakAuthService) IsAnAdminToken(ctx context.Context, token string) bool {
	_, claims, err := k.client.DecodeAccessToken(ctx, token, k.realm)
	if err != nil {
		log.Errorf("Error decoding token: %v", err)
		return false
	}
	scopes := (*claims)[tokenScopes].(string)
	return strings.Contains(scopes, adminScope)
}
