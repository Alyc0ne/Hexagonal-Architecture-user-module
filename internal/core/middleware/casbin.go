package middleware

import (
	"net/http"

	"github.com/LordMoMA/Hexagonal-Architecture/internal/adapters/repository"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

type (
	MiddlewareHandlerService interface {
		CasbinMiddleware() gin.HandlerFunc
	}

	middlewareHandler struct {
		enforcer *casbin.Enforcer
		userRepo repository.UserRepositoryService
	}
)

func NewCasbinMiddleware(enforcer *casbin.Enforcer, userRepo repository.UserRepositoryService) middlewareHandler {
	return middlewareHandler{enforcer, userRepo}
}

func (h *middlewareHandler) CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := h.userRepo.FindUserByEmail(c.Request.Header.Get("email")) // จริงๆต้องใช้ access_token
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		sub := user.Role          // subject (user)
		obj := c.Request.URL.Path // object (URL path)
		act := c.Request.Method   // action (HTTP method)

		allowed := h.enforcer.Enforce(sub, obj, act)
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		c.Next()
	}
}
