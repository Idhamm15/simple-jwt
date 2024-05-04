package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/golang-jwt/jwt/v4"
)

// JWT key: dipakai untuk membuat signature JWT token.
var jwtKey = []byte("secret-key")

// Data key: user, password, role yang bisa digunakan untuk mengakses API
var users = map[string]*User{
    "aditira": {
        Password: "password1",
        Role:     "admin",
    },
    "dito": {
        Password: "password2",
        Role:     "student",
    },
}

type User struct {
    Password string
    Role     string
}

// Struct untuk membaca request body JSON
type Credentials struct {
    Password string `json:"password"`
    Username string `json:"username"`
}

// Struct Claims digunakan sebagai object yang akan di encode atau di parse oleh JWT
// jwt.StandardClaims ditambahkan sebagai embedded type untuk memudahkan proses encoding, parsing dan validasi JWT
type Claims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.StandardClaims
}


func main() {
    fmt.Println("Starting Server at port :8080")
    log.Fatal(http.ListenAndServe(":8080", Routes()))
}

func Routes() *http.ServeMux {
    mux := http.NewServeMux()

    mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        var creds Credentials
        // JSON body diconvert menjadi credential struct & return bad request ketika terjadi kesalahan saat proses decoding json
        err := json.NewDecoder(r.Body).Decode(&creds)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        // Cek apakah username dan password ada dan sesuai dengan yang ada di Data Key & return unauthorized jika password salah
        expectedPassword, ok := users[creds.Username]
        if !ok || expectedPassword.Password != creds.Password {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        // Deklarasi expiry time untuk token jwt
        expirationTime := time.Now().Add(5 * time.Minute)
        // Buat claims berisi data username dan role yang akan kita embed ke JWT
        claims := &Claims{
        Username: creds.Username,
        Role:     expectedPassword.Role,
        StandardClaims: jwt.StandardClaims{
                // expiry time menggunakan time millisecond
                ExpiresAt: expirationTime.Unix(),
            },
        }

        // Buat token menggunakan encoded claim dengan salah satu algoritma yang dipakai
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

        // Buat JWT string dari token yang sudah dibuat menggunakan JWT key yang telah dideklarasikan (proses encoding JWT)
        tokenString, err := token.SignedString(jwtKey)
        if err != nil {
            // return internal error ketika ada kesalahan saat pembuatan JWT string
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        // Set token string ke dalam cookie response
        http.SetCookie(w, &http.Cookie{
            Name:    "token",
            Value:   tokenString,
            Expires: expirationTime,
        })
    })

    // halaman /admin hanya dapat diakses oleh user yang sudah diotentikasi dan memiliki role admin
    mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
        // Ambil token dari cookie yang di kirim ketika request
        c, err := r.Cookie("token")
        if err != nil {
            if err == http.ErrNoCookie {
                // return unauthorized ketika token kosong
                w.WriteHeader(http.StatusUnauthorized)
                return
            }
            // return bad request ketika field token tidak ada
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        // Ambil value dari cookie token
        tknStr := c.Value

        // Deklarasi variable claims yang akan kita isi dengan data hasil parsing JWT
        claims := &Claims{}

        // parse JWT token ke dalam claims
        tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil {
            if err == jwt.ErrSignatureInvalid {
                // return unauthorized ketika ada kesalahan saat parsing token
                w.WriteHeader(http.StatusUnauthorized)
                return
            }
            // return bad request ketika field token tidak ada
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        //return unauthorized ketika token sudah tidak valid (biasanya karena token expired)
        if !tkn.Valid {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        //return unauthorized ketika role user tidak sesuai dengan role admin
        if claims.Role != "admin" {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        // return data dalam claims, yaitu username yang telah didefinisikan di variable claims
        w.Write([]byte(fmt.Sprintf("Welcome Admin %s!", claims.Username)))
    })

    // halaman /profile dapat diakses oleh user baik admin maupun bukan selama terotentikasi.
    mux.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
        // Ambil token dari cookie yang di kirim ketika request
        c, err := r.Cookie("token")
        if err != nil {
            if err == http.ErrNoCookie {
                // return unauthorized ketika token kosong
                w.WriteHeader(http.StatusUnauthorized)
                return
            }
            // return bad request ketika field token tidak ada
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        // Ambil value dari cookie token
        tknStr := c.Value

        // Deklarasi variable claims yang akan kita isi dengan data hasil parsing JWT
        claims := &Claims{}

        //parse JWT token ke dalam claims
        tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil {
            if err == jwt.ErrSignatureInvalid {
                // return unauthorized ketika ada kesalahan ketika parsing token
                w.WriteHeader(http.StatusUnauthorized)
                return
            }
            // return bad request ketika field token tidak ada
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        //return unauthorized ketika token sudah tidak valid (biasanya karena token expired)
        if !tkn.Valid {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        // return data dalam claims, yaitu username yang telah didefinisikan di variable claims
        w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
    })

    return mux
}