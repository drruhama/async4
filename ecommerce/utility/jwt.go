package utility

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWT struct {
	// untuk user id
	Id int `json:"id"`

	// untuk validasi apakah tokennya masih valid atau
	// sudah expired
	Expires time.Time `json:"expires"`
}

var secretKey string // secret key yang akan kita gunakan
var expired int      // in minute

// fungsi untuk nge init secret key dan expired token
func InitToken(secret string, expiredToken int) {
	secretKey = secret
	expired = expiredToken
}

// untuk membuat object JWT
func NewJWT(id int) JWT {
	return JWT{
		Id: id,
		// set expire time
		Expires: time.Now().Add(time.Duration(expired) * time.Minute),
	}
}

func (j JWT) GenerateToken() (tokString string, err error) {
	// membuat object token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      j.Id,
		"expires": j.Expires,
	})

	// proses sign token menggunakan secret key
	tokString, err = token.SignedString([]byte(secretKey))
	return
}

func VerifyToken(tokString string) (token JWT, err error) {
	// proses parsing token dari string menjadi object token
	jwtToken, err := jwt.Parse(tokString, func(t *jwt.Token) (interface{}, error) {
		// validasi, apakah method untuk generate token tadi menggunakan
		// salah satu method HMAC atau engga
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			// jika engga, maka return invalid method
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		// return secret key untuk proses decode nya
		return []byte(secretKey), nil
	})
	if err != nil {
		return
	}

	// type assertion claims untuk dapatin valuenya
	// dalam bentuk map[string]interface{}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		err = fmt.Errorf("invalid token")
		return
	}

	// get id from claims, hasilnya adalah sebuah interface
	id := claims["id"]

	// get expires time from claims, hasilnya adalah sebuah string
	expires := fmt.Sprintf("%v", claims["expires"])

	// parse expires(string) to time
	expiresTime, err := time.Parse(time.RFC3339, expires)
	if err != nil {
		return
	}

	// validate apakah token masih valid atau sudah expired
	if time.Now().After(expiresTime) {
		err = fmt.Errorf("token expired")
		return
	}

	// ubah id jadi int
	idInt, err := strconv.Atoi(fmt.Sprintf("%v", id))
	if err != nil {
		return
	}

	// buat object JWT dengan memasukkan id
	token = NewJWT(idInt)

	return
}
