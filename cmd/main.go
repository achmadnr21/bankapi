/*
Welcome to BankAPI
This is a simple REST API for managing bank accounts, users, branches, and transaction fees.
It is built using Go and Gin framework.
*/
package main

import (
	"database/sql"
	"fmt"

	"github.com/achmadnr21/bankapi/config"
	pgsql "github.com/achmadnr21/bankapi/infrastructure/postgresql"
	"github.com/achmadnr21/bankapi/internal/handler"
	"github.com/achmadnr21/bankapi/internal/middleware"
	"github.com/achmadnr21/bankapi/internal/repository"
	"github.com/achmadnr21/bankapi/internal/usecase"
	"github.com/achmadnr21/bankapi/internal/utils"
	gin_api "github.com/achmadnr21/bankapi/service"
	"github.com/gin-gonic/gin"
)

func main() {
	var sc config.Config
	sc.LoadConfig()
	err := pgsql.InitDB(sc.DbHost, int32(sc.DbPort), sc.DbUser, sc.DbPassword, sc.DbName, sc.DbSsl)
	if err != nil {
		fmt.Println(err)
	}

	var db *sql.DB = pgsql.GetDB()
	defer db.Close()

	// load the jwt secret
	// and refresh secret from the config
	// and initialize the jwt utils
	utils.JwtInit(sc.JwtSecret, sc.RefreshSecret)

	var api gin_api.RESTapi
	var apiV *gin.RouterGroup = api.Init()

	if apiV == nil {
		fmt.Println("API is nil")
	} else {
		fmt.Println("API is prepared")
	}

	// repositories
	userRepository := repository.NewUserRepository(db)
	userRoleRepository := repository.NewUserRoleRepository(db)
	branchRepository := repository.NewBranchRepository(db)
	accountTypeRepository := repository.NewAccountTypeRepository(db)
	currencyRepository := repository.NewCurrencyRepository(db)
	accountRepository := repository.NewAccountRepository(db)
	transactionFeeRepository := repository.NewTransactionFeeRepository(db)

	// usecases
	authUsecase := usecase.NewAuthUsecase(userRepository)
	userUsecase := usecase.NewUserUsecase(userRepository, userRoleRepository)
	branchUsecase := usecase.NewBranchUsecase(branchRepository, userRoleRepository, userRepository)
	accountTypeUsecase := usecase.NewAccountTypeUsecase(accountTypeRepository, userRepository, userRoleRepository)
	currencyUsecase := usecase.NewCurrencyUsecase(currencyRepository, userRepository, userRoleRepository)
	accountUsecase := usecase.NewAccountUsecase(accountRepository, userRepository, userRoleRepository, branchRepository, accountTypeRepository, currencyRepository)
	transactionFeesUsecase := usecase.NewTransactionFeeUsecase(transactionFeeRepository, userRepository, userRoleRepository)

	// handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	branchHandler := handler.NewBranchHandler(branchUsecase)
	accountTypeHandler := handler.NewAccountTypeHandler(accountTypeUsecase)
	currencyHandler := handler.NewCurrencyHandler(currencyUsecase)
	accountHandler := handler.NewAccountHandler(accountUsecase)
	transactionFeeHandler := handler.NewTransactionFeeHandler(transactionFeesUsecase)
	// middleware

	// === START ROUTES ===
	// Auth Router : Public
	auth := apiV.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Users Router : Private
	users := apiV.Group("/users")
	users.Use(middleware.JWTAuthMiddleware)
	{
		users.GET("/search", userHandler.Search)
		users.POST("", userHandler.AddUser)
		users.GET("/:nik", userHandler.GetByNIK)
	}

	// Branch Router : Private
	branches := apiV.Group("/branches")
	branches.Use(middleware.JWTAuthMiddleware)
	{
		branches.GET("", branchHandler.GetAllBranches)
		branches.GET("/:id", branchHandler.GetBranchByID)
		branches.POST("", branchHandler.AddBranch)
		branches.PUT("/:id", branchHandler.UpdateBranch)
		// branches.DELETE("/:id", branchHandler.DeleteBranch)
	}

	// Account Type Router : Private
	accountTypes := apiV.Group("/account-types")
	accountTypes.Use(middleware.JWTAuthMiddleware)
	{
		accountTypes.GET("", accountTypeHandler.GetAllAccountTypes)
		accountTypes.GET("/:id", accountTypeHandler.GetAccountTypeByID)
		accountTypes.POST("", accountTypeHandler.AddAccountType)
		accountTypes.PUT("/:id", accountTypeHandler.UpdateAccountType)
		// accountTypes.DELETE("/:id", accountTypeHandler.DeleteAccountType)
	}

	// Add Currency Router : Private
	currencies := apiV.Group("/currencies")
	currencies.Use(middleware.JWTAuthMiddleware)
	{
		currencies.GET("", currencyHandler.GetAllCurrencies)
		currencies.GET("/:id", currencyHandler.GetCurrencyByID)
		currencies.POST("", currencyHandler.AddCurrency)
		currencies.PUT("/:id", currencyHandler.UpdateCurrency)
		// currencies.DELETE("/:id", branchHandler.DeleteCurrency)
	}

	// Account Router : Private
	accounts := apiV.Group("/accounts")
	accounts.Use(middleware.JWTAuthMiddleware)
	{
		accounts.GET("", accountHandler.GetAllAccounts)
		accounts.GET("/:accnumber", accountHandler.GetAccountByAccountNumber)
		accounts.POST("/:branch/:acctype/:nik", accountHandler.AddAccount)
		accounts.PUT("/:accnumber", accountHandler.UpdateAccount)
		// accounts.DELETE("/:id", branchHandler.DeleteAccount)
	}
	// Transaction Fees Management Routes
	transactionFees := apiV.Group("/transaction-fees")
	transactionFees.Use(middleware.JWTAuthMiddleware)
	{
		transactionFees.GET("", transactionFeeHandler.GetAllTransactionFees)
		transactionFees.GET("/:id", transactionFeeHandler.GetTransactionFeeByID)
		transactionFees.POST("", transactionFeeHandler.AddTransactionFee)
		transactionFees.PUT("/:id", transactionFeeHandler.UpdateTransactionFee)
		transactionFees.DELETE("/:id", transactionFeeHandler.DeleteTransactionFee)
	}
	// === END ROUTES ===

	service_config := fmt.Sprintf(":%d", sc.ServicePort)
	fmt.Printf("\nService running on port %s \n", service_config)
	api.Router.Run(service_config)

}
