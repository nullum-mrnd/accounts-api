package main

import (
	"accounts/api"
	"accounts/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountAPI struct {
	db *db.Db
}

func NewAccountAPI(d *db.Db) *AccountAPI {
	return &AccountAPI{db: d}
}

func (a AccountAPI) CreateAccount(c *gin.Context) {
	var newAccount api.NewAccount
	if err := c.ShouldBindBodyWithJSON(&newAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := a.db.CreateAccount(*newAccount.Username, *newAccount.Phone)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, acc)
}

func (a AccountAPI) DeleteAccount(c *gin.Context, id string) {
	intId, _ := strconv.Atoi(id)
	_, err := a.db.GetAccountByID(intId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err = a.db.DeleteAccount(intId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)

}

func (a AccountAPI) GetAccounts(c *gin.Context) {
	accs, err := a.db.GetAccounts()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	apiAccs := make([]api.Account, len(accs))
	for i, acc := range accs {
		apiAccs[i] = api.Account{
			Id:       &acc.UserID,
			Username: &acc.Username,
			Phone:    &acc.Phone,
		}
	}

	c.JSON(http.StatusOK, apiAccs)
}

func (a AccountAPI) GetAccountById(c *gin.Context, id string) {
	intId, _ := strconv.Atoi(id)

	acc, err := a.db.GetAccountByID(intId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.Account{
		Id:       &acc.UserID,
		Username: &acc.Username,
		Phone:    &acc.Phone,
	})
}

func (a AccountAPI) UpdateAccount(c *gin.Context, id string) {
	intId, _ := strconv.Atoi(id)

	acc, err := a.db.GetAccountByID(intId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	//var newAccount api.NewAccount
	if err := c.ShouldBindBodyWithJSON(&acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = a.db.UpdateAccount(intId, *acc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//(*acc).Username = *newAccount.Username
	//(*acc).Phone = *newAccount.Phone
	c.JSON(http.StatusOK, api.Account{
		Id:       &acc.UserID,
		Username: &acc.Username,
		Phone:    &acc.Phone,
	})
}
