package controller

import (
	"controller/helpers"
	"controller/maintypes"
	"controller/model"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Createuser(c *gin.Context) {

	var user model.LoUser
	var count int
	if err := c.BindJSON(&user); err != nil {
		// PtrToLogger.Error("sgt_portal_controller", zap.String("message", fmt.Sprintf("Create user: Unable to bind payload failed because %s", err.Error())), zap.String("sendto", string(sgttypes.Remote)))
		c.JSON(http.StatusBadRequest, gin.H{"message": "SG Cloud is unable to process your request, please report code 400 to SGT Support"})
		return
	}
	checkExistingUser := helpers.MariaDB.Table("users").Select("count(*)").Where("email = ?", user.Email).Find(&count)

	if checkExistingUser.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot create User please, Please report to Support"})
		return
	}
	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User already exist, please try to login"})
		return
	}

	// jwt, err := helpers.Portalusertoken(user.Username, user)

	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "SG Cloud is unable to process your request, please report code 400 to SGT Support"})
	// 	return
	// }
	encwrd := helpers.GetMD5Hash(user.Password)
	user.Password = encwrd

	newrecord := model.Users{Username: user.Username, Email: user.Email, Password: user.Password, Isadmin: "N"}
	result := helpers.MariaDB.Create(&newrecord)
	if result.Error != nil {
		panic(result.Error)
	}
	c.JSON(http.StatusOK, gin.H{"Message": "User created successfully"})

}

func Userlogin(c *gin.Context) {

	type Requser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var requestUser Requser
	var user model.LoUser
	// var password string
	if err := c.BindJSON(&requestUser); err != nil {
		// PtrToLogger.Error("sgt_portal_controller", zap.String("message", fmt.Sprintf("Create user: Unable to bind payload failed because %s", err.Error())), zap.String("sendto", string(sgttypes.Remote)))
		c.JSON(http.StatusBadRequest, gin.H{"message": "SG Cloud is unable to process your request, please report code 400 to SGT Support"})
		return
	}

	checkExistingUser := helpers.MariaDB.Table("users").Select("email,username,password").Where("email = ?", requestUser.Email).Find(&user)

	if checkExistingUser.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusForbidden, "message": "Cannot find User, Please, Sign up first to Login"})
		return
	}

	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Cannot find User, Please, Sign up first to Login"})
		return
	}
	encwrd := helpers.GetMD5Hash(requestUser.Password)

	if user.Password != encwrd {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "User provided password is incorrect, Please report to Support"})
		return
	}

	jwt, err := helpers.Portalusertoken(user.Username, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "SG Cloud is unable to process your request, please report code 400 to SGT Support"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user.Email, "jwt_token": jwt})

}

func AuthUser() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var user model.LoUser
		tokenString := ctx.GetHeader("Authorization")
		requestuser := ctx.GetHeader("PortalUser")
		if strings.TrimSpace(tokenString) == "" || strings.TrimSpace(requestuser) == "" {
			ctx.JSON(401, gin.H{"message": "Unauthorized Request, Please login again!"})
			ctx.Abort()
			return
		}
		checkExistingUser := helpers.MariaDB.Table("users").Select("email,username,password").Where("email = ?", requestuser).Find(&user)

		if checkExistingUser.Error != nil {
			ctx.JSON(401, gin.H{"message": "Unauthorized Request, Please login again!"})
			// c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot create User please, Please report to Support"})
			return
		}
		jw := helpers.JwtWrapper{
			SecretKey:       fmt.Sprintf("%s*%s_%s", user.Username, user.Password, user.Email),
			Issuer:          strconv.Itoa(maintypes.RandNumber),
			ExpirationHours: 1,
		}
		_, err := jw.ValidateJWTToken(tokenString)
		if err != nil {
			ctx.JSON(401, gin.H{"message": "Unauthorized Request, Please login again!"})
			ctx.Abort()
			return
		}
	}
}

func Getproduct(c *gin.Context) {

	type product struct {
		PID int `json:"pro_id"`
	}

	var pid product
	var productdetail model.Product
	if err := c.BindJSON(&pid); err != nil {
		// PtrToLogger.Error("sgt_portal_controller", zap.String("message", fmt.Sprintf("Create user: Unable to bind payload failed because %s", err.Error())), zap.String("sendto", string(sgttypes.Remote)))
		c.JSON(http.StatusBadRequest, gin.H{"message": "SG Cloud is unable to process your request, please report code 400 to SGT Support"})
		return
	}

	checkExistingUser := helpers.MariaDB.Table("product").Select("*").Where("id = ?", pid.PID).Find(&productdetail)

	if checkExistingUser.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot create User please, Please report to Support"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Success", "Product": productdetail})

}

func Addcart(c *gin.Context) {
	pUser := strings.TrimSpace(c.GetHeader("PortalUser"))
	var productdetail model.Product
	var users model.Users
	if len(pUser) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User credentials invalid"})
		return
	}
	Pid := c.Param("productid")
	if len(Pid) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User credentials invalid"})
		return
	}
	checkExistingUser := helpers.MariaDB.Table("users").Select("id,email,username").Where("email = ?", pUser).Find(&users)

	if checkExistingUser.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusForbidden, "message": "Cannot find User, Please, Sign up first to Login"})
		return
	}
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	checkProduct := helpers.MariaDB.Table("product").Select("*").Where("id = ?", Pid).Find(&productdetail)
	if checkProduct.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot create User please, Please report to Support"})
		return
	}
	newrecord := model.Cartdetails{Productid: productdetail.Id, Customerid: users.Id, Eventtimestamp: formattedTime}

	result := helpers.MariaDB.Create(&newrecord)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot add to cart, Please report to Support"})
		return
	}
	// fmt.Println(formattedTime)
	c.JSON(http.StatusOK, gin.H{"Message": "Added to cart successfully"})

}

func Getcart(c *gin.Context) {
	pUser := strings.TrimSpace(c.GetHeader("PortalUser"))

	type User struct {
		Id int
	}
	var count int
	var users User
	if len(pUser) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User credentials invalid"})
		return
	}

	checkExistingUser := helpers.MariaDB.Table("users").Select("id").Where("email = ?", pUser).Find(&users)

	if checkExistingUser.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusForbidden, "message": "Cannot find User, Please, Sign up first to Login"})
		return
	}

	productincart := helpers.MariaDB.Table("cartdetails").Select("count(*)").Where("customerid = ?", users.Id).Find(&count)

	if productincart.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusForbidden, "message": "Cannot find product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"prodcutcount": count})

}
