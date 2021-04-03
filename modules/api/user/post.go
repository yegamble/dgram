package user

import (
	"dgram/database"
	keyUtil "dgram/modules/util"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID       uuid.UUID `json:"id" gorm:"primary_key"`
	UserID   uuid.UUID `json:"user_id" gorm:"text"`
	Text     string    `json:"text" gorm:"type:text"`
	Images   []string  `json:"images" gorm:"type:text"`
	Videos   []string  `json:"videos" gorm:"type:text"`
	Comments []string  `json:"comments" gorm:"type:text"`
	Votes    []Vote    `json:"votes" gorm:"type:text"`
	gorm.Model
}

type Comment struct {
	ID    uuid.UUID `json:"id" gorm:"primary_key"`
	Text  string    `json:"text" gorm:"text"`
	Votes []Vote    `json:"votes" gorm:"type:text"`
}

type Vote struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func CreateNewPost(c *fiber.Ctx) error {

	db := database.DBConn

	UploadPostMedia(c)

	var NewPost Post
	user := FindUser(c.Params("id"))

	NewPost.ID = uuid.New()
	NewPost.UserID = user.ID

	error := c.BodyParser(&NewPost)
	if error != nil || isValidPost(&NewPost) {
		return FailedTransaction(c)
	}

	db.Save(&NewPost)
	return c.Status(fiber.StatusOK).JSON(&NewPost)
}

func UploadPostMedia(c *fiber.Ctx) (string, error) {

	// Check for errors:

	file, err := c.FormFile("file")
	if err != nil {
		return "", nil
	}

	dir := fmt.Sprintf("./uploads/%s", file.Filename)
	c.SaveFile(file, dir)

	hash, err := keyUtil.UploadToIPFS(dir)
	if err != nil {
		return "", nil
	}

	return hash, nil

}

func findPost() {

}

func SavePost(p *Post) {

}

func isValidPost(p *Post) bool {

	if len(p.Images) == 0 && len(p.Videos) == 0 {
		return false
	}

	return true
}
