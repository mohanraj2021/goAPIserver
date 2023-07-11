package helpers

import (
	"controller/maintypes"
	"controller/model"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// var PtrToLogger *zap.Logger
var MariaDBConnStr string

func Wakeup() bool {

	if err := godotenv.Load(".env"); err != nil {
		// lg.Error("sgt_portal_controller", zap.String("message", fmt.Sprintf("env file could not be loaded because %s", err.Error())), zap.String("sendto", string(maintypes.Local)))
		return false
	}

	MariaDBConnStr = fmt.Sprintf("%s:%s@tcp(%v:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), "college_proj_db")

	return true

}

func Portalusertoken(username string, l_user model.LoUser) (sgntoken string, err error) {
	j := JwtWrapper{
		SecretKey:       fmt.Sprintf("%s*%s_%s", l_user.Username, l_user.Password, l_user.Email),
		Issuer:          strconv.Itoa(maintypes.RandNumber),
		ExpirationHours: 1,
	}
	sgntoken, err = j.GenerateJWTToken(username)
	return sgntoken, err
}

func Ecrypt() (encpwd string, encerr error) {
	key := []byte("0123456789abcdef0123456789abcdef")

	// The data to be encrypted
	plaintext := []byte("Hello, World!")

	// Create a new AES cipher block using the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Create a new byte array for the ciphertext
	// paddedPlaintext := pkcs7.Pad(plaintext, aes.BlockSize)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// Generate a random initialization vector (IV)
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Encrypt the data
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// Print the ciphertext
	encWrd := fmt.Sprintf("Ciphertext: %x\n", ciphertext)

	// Decrypt the data
	mode = cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], ciphertext[aes.BlockSize:])

	// Remove padding from the plaintext
	plaintext = PKCS7Unpad(ciphertext[aes.BlockSize:])

	// Print the decrypted plaintext
	fmt.Println("Plaintext:", string(plaintext))
	return encWrd, nil
}

// PKCS7Unpad removes padding from the given ciphertext using the PKCS#7 scheme
func PKCS7Unpad(ciphertext []byte) []byte {
	padding := ciphertext[len(ciphertext)-1]
	return ciphertext[:len(ciphertext)-int(padding)]
}

func StringRepair(text string) string {
	res1 := strings.ToLower(text)
	trimmed := strings.ReplaceAll(res1, " ", "")
	return GetMD5Hash(trimmed)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
