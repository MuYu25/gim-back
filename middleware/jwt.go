package middleware

import (
	"errors"
	"net/http"
	"project/utils"
	"project/utils/errmsg"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	JwtKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(utils.JwtKey),
	}
}

type MyClaims struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 定义错误
var (
	ErrTokenExpired     = errors.New("token已过期,请重新登录。")
	ErrTokenNotValidYet = errors.New("token无效,请重新登录。")
	ErrTokenMalformed   = errors.New("token不正确,请重新登录。")
	ErrTokenInvalid     = errors.New("这不是一个token,请重新登录。")
)

// CreateToken 生成token
func (j *JWT) CreateToken(claims MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JwtKey)
}

// ParserToken 解析token
func (j *JWT) ParserToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.JwtKey, nil
	})
	return token, err
}

// verifyToken 验证token
func (j *JWT) verifyToken(tokenString string) error {
	token, err := j.ParserToken(tokenString)
	// 验证token
	if token.Valid {
		return nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return ErrTokenMalformed
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return ErrTokenExpired
	} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return ErrTokenInvalid
	} else {
		return ErrTokenNotValidYet
	}
}

// JwtToken jwt中间件
// todo 优化此类代码
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			code = errmsg.ERROR_TOKEN_EXIST
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		j := NewJWT()
		// 解析token
		err := j.verifyToken(tokenHeader)
		if err != nil {
			if errors.Is(err, ErrTokenExpired) {
				c.JSON(http.StatusOK, gin.H{
					"status":  errmsg.ERROR,
					"message": "token授权已过期,请重新登录",
					"data":    nil,
				})
				c.Abort()
				return
			}
			// 其他错误
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR,
				"message": err.Error(),
				"data":    nil,
			})
			c.Abort()
			return
		}

		//c.Set("username",)
		c.Next()
	}
}
