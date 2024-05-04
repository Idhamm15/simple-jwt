Untuk mencoba kode di atas menggunakan Postman, kita dapat mengikuti langkah-langkah berikut:

- Jalankan server dengan menjalankan perintah go run nama_file.go pada terminal.
- Buka Postman dan buat request baru dengan metode POST. Isikan alamat URL dengan http://localhost:8080/login
- Pada tab body, pilih raw dan pilih JSON(application/json) sebagai format yang digunakan.
- Isikan body dengan format JSON berikut {"username":"aditira", "password":"password1"}
- Klik tombol send untuk mengirimkan request.
- Jika login berhasil, server akan mengeluarkan token JWT yang diterima dalam bentuk cookie. Anda dapat menggunakan token ini untuk mengirimkan permintaan ke halaman yang dibatasi oleh otorisasi.
- Buat request baru dengan metode GET, isikan alamat URL dengan http://localhost:8080/admin
- Klik tombol send untuk mengirimkan request.
- Jika token yang digunakan valid, server akan mengirimkan respon yang sesuai. Jika tidak, server akan mengembalikan kode respon 401 Unauthorized.