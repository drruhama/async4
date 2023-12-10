package routerChi

import (
	"context"
	"ecommerce/utility"
	"log"
	"net/http"
	"strings"
	"time"
)

// disini, kita membuat 1 function untuk middleware Tracer
//func Tracer(h http.Handler) http.Handler {
// lalu function ini akan nge-return sebuah handlerFunc
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		log.Println("Before Request")

// proses ini akan meneruskan request ke next handler
//		h.ServeHTTP(w, r)

//		log.Println("After Request")
//	})
//}

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get current time
		now := time.Now()

		// setup log for incoming request
		log.Printf("method=%v path=%v type=[INFO] message='incoming request'", r.Method, r.URL.Path)

		// pass to next handler
		h.ServeHTTP(w, r)

		// compare process and get time by ms
		end := time.Since(now).Milliseconds()

		// log for after request with response time
		log.Printf("method=%v path=%v type=[INFO] message='finish request' response_time=%vms", r.Method, r.URL.Path, end)
	})
}

func CheckToken(h http.Handler) http.Handler {
	// untu flownya sama, namun kita hanya perlu
	// melakukan pengecekan before requestnya aja.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// coding pengecekannya akan ada disini
		bearer := r.Header.Get("Authorization")

		// jika gaada bearer token, maka return unauthorized
		if bearer == "" {
			resp := APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "UNAUTHORIZED",
				Error:   "no token provided",
			}
			WriteJsonResponse(w, resp)
			return
		}

		// lakukan split string
		// artinya stringnya akan kita potong berdasarkan "Bearer "
		tokSlice := strings.Split(bearer, "Bearer ")

		// expectnya akan ada 2 data
		if len(tokSlice) < 2 {
			resp := APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "UNAUTHORIZED",
				Error:   "invalid token",
			}
			WriteJsonResponse(w, resp)
			return
		}

		// data pada index 0 nya akan kosong
		// dan pada index 1 adalah isi dari tokennya (token string)
		tokString := tokSlice[1]

		// setelah itu, lakukan verifikasi token tersebut
		token, err := utility.VerifyToken(tokString)

		// jika error, return unauthorized
		if err != nil {
			resp := APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "UNAUTHORIZED",
				Error:   err.Error(),
			}
			WriteJsonResponse(w, resp)
			return
		}

		// nah bagiain ini, kita akan mengirim value dari id nya
		// via context.
		ctx := context.WithValue(r.Context(), "AUTH_ID", token.Id)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}
