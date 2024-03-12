package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type KeycloakConfiguration struct {
	KeycloakClient       string
	KeycloakClientSecret string
	KeycloakRealm        string
	KeycloakUrl          string
}

func getKeycloakConfig(env *viper.Viper, securityEnabled bool) *KeycloakConfiguration {
	if !securityEnabled {
		log.Warnf("Security disabled - Development environment only")
		return nil
	}

	client := getEnvString(env, "keycloak.client")
	secret := getEnvString(env, "keycloak.secret")
	realm := getEnvString(env, "keycloak.realm")
	url := getEnvString(env, "keycloak.url")
	log.Info("Security enabled - Production environment")
	return &KeycloakConfiguration{client, secret, realm, url}
}
