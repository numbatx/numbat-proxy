package api

import (
	"fmt"
	"reflect"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/numbatx/numbat-proxy/api/address"
	"github.com/numbatx/numbat-proxy/api/transaction"
	"gopkg.in/go-playground/validator.v8"
)

type validatorInput struct {
	Name      string
	Validator validator.Func
}

// Start will boot up the api and appropriate routes, handlers and validators
func Start(numbatProxyFacade NumbatProxyHandler, port int) error {
	ws := gin.Default()
	ws.Use(cors.Default())

	err := registerValidators()
	if err != nil {
		return err
	}
	registerRoutes(ws, numbatProxyFacade)

	return ws.Run(fmt.Sprintf(":%d", port))
}

func registerRoutes(ws *gin.Engine, numbatProxyFacade NumbatProxyHandler) {
	addressRoutes := ws.Group("/address")
	addressRoutes.Use(WithNumbatProxyFacade(numbatProxyFacade))
	address.Routes(addressRoutes)

	txRoutes := ws.Group("/transaction")
	txRoutes.Use(WithNumbatProxyFacade(numbatProxyFacade))
	transaction.Routes(txRoutes)
}

func registerValidators() error {
	validators := []validatorInput{
		{Name: "skValidator", Validator: skValidator},
	}
	for _, validatorFunc := range validators {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			err := v.RegisterValidation(validatorFunc.Name, validatorFunc.Validator)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// skValidator validates a secret key from user input for correctness
func skValidator(
	_ *validator.Validate,
	_ reflect.Value,
	_ reflect.Value,
	_ reflect.Value,
	_ reflect.Type,
	_ reflect.Kind,
	_ string,
) bool {
	return true
}

// WithNumbatProxyFacade middleware will set up an NumbatFacade object in the gin context
func WithNumbatProxyFacade(numbatProxyFacade NumbatProxyHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("numbatProxyFacade", numbatProxyFacade)
		c.Next()
	}
}
