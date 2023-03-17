
# Authify

Authify is authentication starter pack using Golang and Mysql. Authify dessigned to be fast and independent. To provide stable and fast HTTP Server, I use Echo as a framework.


## Features

- Login
- Register
- Forget Password
- Reset Password
- CRUD User
- Simple Profile with Level Middleware
- JWT Authentication
- Email Confirmation
- Queue Email using MySQL


## Usage

After clone this project you will find this on `main.go`

```bash
func main() {
    go queue.EmailQueue()
	initApp()
	//migrateAndSeed()
}
```

You can run `migrateAndSeed` to run the migration and seed dummy data by uncommand the code.

***note : it's important to run migration before you run the project**

To run this project you only need type `go run main.go` on terminal

## Config
To change the MySQL Config, SMTP, and other things, you can open `config` folder. Sometime i make the config by hardcode them.
## Support

Feel free to answer your quetions 👋

- [Instagram](https://instagram.com/mufthi_ryanda)
- [Email](mailto:mufthi.ryan@gmail.com)

