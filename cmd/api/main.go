package main

import (
	"log"

	"authservice/internal/authorization/jwt"
	"authservice/internal/authservice"
	"authservice/internal/config"
	"authservice/internal/store/redis"
	"authservice/internal/transport"
)

func main() {

	config := config.LoadConfigFromEnv()

	client := redis.NewClient(&config.RedisConfig)

	jwtValidator := jwt.NewValidator(
		jwt.ValidatorConfig{
			Secret: []byte(config.AuthServiceConfig.Secret),
		},
	)

	authService := authservice.NewAuthService(
		jwtValidator,
		client,
	)

	if err := transport.NewServer(authService, &config.AuthServiceConfig).ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
