package cookie

import (
	"github.com/gin-gonic/gin"
)

type CookieConfig struct {
	Domain   string
	Secure   bool
	HttpOnly bool
	Path     string
}

func NewCookieConfig(domain string, secure bool) *CookieConfig {
	return &CookieConfig{
		Domain:   domain,
		Secure:   secure,
		HttpOnly: true,
		Path:     "/",
	}
}

func (c *CookieConfig) SetAccessToken(ctx *gin.Context, token string) {
	ctx.SetCookie(
		"access_token",
		token,
		900,
		c.Path,
		c.Domain,
		c.Secure,
		c.HttpOnly,
	)
}

func (c *CookieConfig) SetRefreshToken(ctx *gin.Context, token string) {
	ctx.SetCookie(
		"refresh_token",
		token,
		604800,
		c.Path,
		c.Domain,
		c.Secure,
		c.HttpOnly,
	)
}

func (c *CookieConfig) ClearAuthCookies(ctx *gin.Context) {
	ctx.SetCookie(
		"access_token", 
		"", 
		-1, 
		c.Path,
		c.Domain,
		c.Secure,
		c.HttpOnly,
	)
}

func (c *CookieConfig) GetAccessToken(ctx *gin.Context) (string, error) {
	return ctx.Cookie("access_token")
}


func (c *CookieConfig) GetRefreshToken(ctx *gin.Context) (string, error) {
	return ctx.Cookie("refresh_token")
}
