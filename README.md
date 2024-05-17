# Online Shope Project
1. Jalankan docker postgresql
```
docker run --name postgresql -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=toko -d -p 5432:5432  postgres:16
```

2. Export ENV yang dibutuhkan
```
export DB_URI=postgres://user:password@localhost:5432/toko?sslmode=disable
export ADMIN_SECRET=secret
```

3. Cara menjalankan program
```
go run .

```