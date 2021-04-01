package user

import (
	"dgram/modules/api/wallet"
	keyUtil "dgram/modules/util"
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	Uuid            uuid.UUID  `json:"id" gorm:"primary_key"`
	FirstName       *string    `json:"first_name"`
	LastName        *string    `json:"last_name"`
	Email           *string    `json:"email"`
	Username        string     `json:"username"`
	DateOfBirth     *time.Time `json:"date_of_birth"`
	Gender          *string    `json:"gender"`
	CurrentCity     *string    `json:"current_city"`
	HomeTown        *string    `json:"hometown"`
	Bio             string     `json:"bio"`
	ProfilePhoto    string     `json:"profile_photo"`
	Password        string     `json:"password"`
	Wallet          string     `json:"wallet"`
	Posts           []Post     `json:"posts"`
	Friends         []User     `json:"friends"`
	PGPKey          string     `json:"pgp_key"`
	DateTimeEdited  time.Time  `json:"datetime_modified"`
	DateTimeCreated time.Time  `json:"datetime_created"`
}

type HashConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

type Users []User

func IndexPage(w http.ResponseWriter, r *http.Request) {}

//Function creates a new User
func CreateNewUser(ctx *fiber.Ctx) error {

	//var body User
	var body User
	body.Uuid = uuid.New()
	body.Wallet = wallet.GenerateNewWallet()
	body.Password, _ = encodeToArgon(body.Password)

	err := ctx.BodyParser(&body)
	if err != nil || isValidUser(&body) != nil {
		return faliedTransaction(ctx)
	}

	body.Username, _ = generateUsername(*body.FirstName, *body.LastName)
	body.PGPKey = keyUtil.Fingerprint(body.PGPKey)

	return ctx.Status(fiber.StatusOK).JSON(body)
}

func generateUsername(FirstName string, LastName string) (string, error) {

	if FirstName == "" || LastName == "" {
		return "", errors.New("Name is Invalid")
	}

	format := "%s.%s.%d"
	return fmt.Sprintf(format, FirstName, LastName, rand.Intn(99999+1)), nil
}

//encodes a string input to argon hash
func encodeToArgon(input string) (string, error) {

	c := &HashConfig{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
	}

	// Generate a Salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(input), salt, c.time, c.memory, c.threads, c.keyLen)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash)
	return full, nil

}

//checks if user is valid before saving to database
func isValidUser(user *User) error {
	if user.FirstName == nil || user.LastName == nil {
		return errors.New("empty name")
	}
	return nil
}

func faliedTransaction(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "Server Error Cannot Create User",
	})
}

func successfulTransaction(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}
