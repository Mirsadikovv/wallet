//go:build tools
// +build tools

package tools

import (
	_ "github.com/gin-gonic/gin"
	_ "github.com/uptrace/bun"
	_ "github.com/uptrace/bun/driver/pgdriver"
	_ "github.com/uptrace/bun/extra/bundebug"
	_ "github.com/redis/go-redis/v9"
	_ "github.com/golang-jwt/jwt/v5"
)
