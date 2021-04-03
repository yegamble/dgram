package user

import (
	"dgram/database"
	"dgram/modules/api/wallet"
	keyUtil "dgram/modules/util"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"strings"
	"time"
)

type User struct {
	ID                  uuid.UUID  `json:"id" gorm:"primary_key"`
	FirstName           *string    `json:"first_name" gorm:"type:text"`
	LastName            *string    `json:"last_name" gorm:"type:text"`
	Email               *string    `json:"email" gorm:"type:text"`
	Username            string     `json:"username" gorm:"unique"`
	DateOfBirth         *time.Time `json:"date_of_birth"`
	Gender              *string    `json:"gender" gorm:"type:datetime"`
	CurrentCity         *string    `json:"current_city" gorm:"type:text"`
	HomeTown            *string    `json:"hometown" gorm:"type:text"`
	Bio                 *string    `json:"bio" gorm:"type:text"`
	ProfilePhoto        string     `json:"profile_photo" gorm:"type:text"`
	HeaderPhoto         string     `json:"header_photo" gorm:"type:text"`
	Password            string     `json:"password" gorm:"type:text"`
	Posts               []Post     `json:"posts" gorm:"foreignKey:user_id;references:ID"`
	PrivateKey          string     `json:"private_key" gorm:"type:text"`
	PublicKey           string     `json:"public_key" gorm:"type:text"`
	LastTailTransaction string     `json:"last_tail_transaction" gorm:"type:text"`
	Friends             []User     `json:"friends" gorm:"type:text"`
	PGPKey              string     `json:"pgp_key" gorm:"type:text"`
	gorm.Model
}

type HashConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

func init() {

}

func GetUsers(c *fiber.Ctx) error {
	db := database.DBConn
	var users []User

	//db.Preload(clause.Associations).Find(&users) //asociations
	db.Find(&users)
	return c.Status(fiber.StatusOK).JSON(&users)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user := FindUser(id)
	if user.FirstName == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "false",
			"message": "Profile Not Found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func SaveUser(user *User) {
	db := database.DBConn

	db.Save(&user)
}

func FindUser(id string) User {
	db := database.DBConn

	var user User
	db.First(&user, "id = ?", id)

	return user
}

func UpdateUser(c *fiber.Ctx) error {

	db := database.DBConn
	id := c.Params("id")

	var body User
	db.First(&body, "id = ?", id)

	err := UploadProfilePhoto(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Error Uploading Photo")
	}

	error := c.BodyParser(&body)
	if error != nil || isValidUser(&body) != nil || body.ID.String() != id {
		return FailedTransaction(c)
	}

	body.PGPKey = keyUtil.Fingerprint(body.PGPKey)

	body = FindUser(id)

	err = userFormMapper(&body, c)
	if err != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	data, err := json.Marshal(body)
	_, body.LastTailTransaction, err = wallet.UpdateProfileAddress(body.PublicKey, body.PrivateKey, string(data))
	if err != nil || isValidUser(&body) != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	} else {
		db.Save(&body)
	}

	transaction, err := wallet.GetTransactionJSON(body.LastTailTransaction)
	if err != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	} else if transaction == "" {
		log.Fatal(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Decentralised Transaction not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(&body)
}

func userFormMapper(user *User, c *fiber.Ctx) error {
	*user.FirstName = c.FormValue("first_name")
	*user.LastName = c.FormValue("last_name", *user.LastName)
	*user.Email = c.FormValue("email", *user.Email)
	*user.Gender = c.FormValue("gender", *user.Gender)
	*user.CurrentCity = c.FormValue("current_city", *user.CurrentCity)
	*user.HomeTown = c.FormValue("hometown", *user.HomeTown)
	*user.Bio = c.FormValue("bio", *user.Bio)

	//timeLayout := time.RFC3339
	//t2, err := time.Parse(timeLayout, c.FormValue("date_of_birth"))
	//
	//if err != nil {
	//	fmt.Println("RFC format doesn't work") // You shouldn't see this at all
	//	return err
	//} else {
	//	fmt.Println(t2)
	//}

	return nil
}

func DeleteUser(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	var body User
	db.First(&body, "id = ?", id)
	if body.FirstName == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "false",
			"message": "Profile Not Found",
		})
	}

	db.Delete(&body)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "true",
		"message": "Profile Deleted",
	})
}

//Function creates a new User
func CreateNewUser(ctx *fiber.Ctx) error {

	db := database.DBConn

	//var body User
	var body User
	body.ID = uuid.New()
	body.PrivateKey, body.PublicKey, body.LastTailTransaction = wallet.GenerateNewWallet()
	body.Password = encodeToArgon(body.Password)
	UploadProfilePhoto(ctx)

	err := ctx.BodyParser(&body)
	if err != nil || isValidUser(&body) != nil {
		return FailedTransaction(ctx)
	}

	body.Username, _ = generateUsername(*body.FirstName, *body.LastName)
	body.PGPKey = keyUtil.Fingerprint(body.PGPKey)

	data, err := json.Marshal(body)
	_, body.LastTailTransaction, err = wallet.UpdateProfileAddress(body.PublicKey, body.PrivateKey, string(data))
	if err != nil {
		log.Fatal(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	} else {
		db.Create(&body)
	}

	transaction, err := wallet.GetTransactionJSON(body.LastTailTransaction)
	if err != nil {
		log.Fatal(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	} else if transaction == "" {
		log.Fatal(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Decentralised Transaction not Found",
		})
	}

	var transactionResultObj User
	json.Unmarshal([]byte(transaction), &transactionResultObj)

	return ctx.Status(fiber.StatusOK).JSON(transactionResultObj)
}

func UploadProfilePhoto(c *fiber.Ctx) error {
	// Get first file from form field "profile_photo":
	file, err := c.FormFile("profile_photo")
	if err != nil || file == nil {
		return nil
	}

	id := c.Params("id")

	// Check for errors:
	if err == nil {

		dir := fmt.Sprintf("./uploads/%s", file.Filename)
		c.SaveFile(file, dir)

		hash, err := keyUtil.UploadToIPFS(dir)
		if err != nil {
			return nil
		}

		user := FindUser(id)
		user.ProfilePhoto = hash
		SaveUser(&user)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":    true,
			"ipfs_hash": hash,
		})
	} else {
		return err
	}
}

func generateUsername(FirstName string, LastName string) (string, error) {
	db := database.DBConn
	if FirstName == "" || LastName == "" {
		return "", errors.New("Name is Invalid")
	}

	format := "%s.%s.%d"
	result := strings.ToLower(fmt.Sprintf(format, FirstName, LastName, rand.Intn(99999+1)))

	var body User
	queryError := db.First(&body, "username = ?", result).Error

	if queryError == nil {
		var err error
		result, err = generateUsername(FirstName, LastName)
		if err != nil {
			return "", nil
		}
	}

	return result, nil
}

//encodes a string input to argon hash
func encodeToArgon(input string) string {

	c := &HashConfig{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
	}

	// Generate a Salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return ""
	}

	hash := argon2.IDKey([]byte(input), salt, c.time, c.memory, c.threads, c.keyLen)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash)
	return full

}

//checks if user is valid before saving to database
func isValidUser(user *User) error {
	if user.FirstName == nil || user.LastName == nil {
		return errors.New("empty name")
	}
	return nil
}

func FailedTransaction(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "Server Error Cannot Make Update",
	})
}

func successfulTransaction(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}
