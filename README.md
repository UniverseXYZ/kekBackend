# barnbridge-backend

This repository will host the BarnBridge backend service/s.

At this point, it has the following features: 
- connection to a postgres database
- sample API with a test endpoint using [Gin](https://github.com/gin-gonic/gin)
- automatically run the migrations (enabled via feature flags)
- uses [Cobra](https://github.com/spf13/cobra) for managing commands
- uses [Viper](https://github.com/spf13/viper) for managing configuration