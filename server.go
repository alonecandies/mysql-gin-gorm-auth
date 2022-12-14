package main

import (
	"github.com/alonecandies/mysql-gin-gorm-auth/api/configs/db"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/controllers"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/middlewares"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/repositories"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	conn           *gorm.DB                    = db.DBConnection()
	userRepository repositories.UserRepository = repositories.NewUserRepository(conn)
	bookRepository repositories.BookRepository = repositories.NewBookRepository(conn)
	jwtService     services.JWTService         = services.NewJWTService()
	authService    services.AuthService        = services.NewAuthService(userRepository)
	userService    services.UserService        = services.NewUserService(userRepository)
	bookService    services.BookService        = services.NewBookService(bookRepository)
	authController controllers.AuthController  = controllers.NewAuthController(authService, jwtService)
	userController controllers.UserController  = controllers.NewUserController(userService, jwtService)
	bookController controllers.BookController  = controllers.NewBookController(bookService, jwtService)
)

func main() {
	defer db.DBClose(conn)
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middlewares.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
		userRoutes.GET("/books", userController.MyBooks)
	}
	
	bookRoutes:= r.Group("api/book", middlewares.AuthorizeJWT(jwtService))
	{
		bookRoutes.GET("/", bookController.AllBooks)
		bookRoutes.GET("/:id", bookController.FindBookById)
		bookRoutes.POST("/", bookController.InsertBook)
		bookRoutes.PUT("/:id", bookController.UpdateBook)
		bookRoutes.DELETE("/:id", bookController.DeleteBook)
	}

	r.Run()
}
