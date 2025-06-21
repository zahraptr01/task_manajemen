package middleware

import (
	"fmt"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userCookie, err := r.Cookie("user_id")
		if err != nil || userCookie.Value == "" {
			// Debug: print cookies yang masuk
			for _, c := range r.Cookies() {
				fmt.Printf("Cookie: %s = %s\n", c.Name, c.Value)
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Debug: tampilkan user_id yang diterima
		fmt.Println("User ID ditemukan:", userCookie.Value)

		next.ServeHTTP(w, r)
	})
}
