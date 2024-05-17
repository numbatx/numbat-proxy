package transaction

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/numbatx/numbat-proxy/api/errors"
	"github.com/numbatx/numbat-proxy/data"
)

// Routes defines transaction related routes
func Routes(router *gin.RouterGroup) {
	router.POST("/send", SendTransaction)
}

// SendTransaction will receive a transaction from the client and propagate it for processing
func SendTransaction(c *gin.Context) {
	ef, ok := c.MustGet("numbatProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInvalidAppContext.Error()})
		return
	}

	var gtx = data.Transaction{}
	err := c.ShouldBindJSON(&gtx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error())})
		return
	}

	signature, err := hex.DecodeString(gtx.Signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrInvalidSignatureHex.Error(), err.Error())})
		return
	}

	txHash, err := ef.SendTransaction(gtx.Nonce, gtx.Sender, gtx.Receiver, gtx.Value, gtx.Data, signature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrTxGenerationFailed.Error(), err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"txHash": txHash})
}
