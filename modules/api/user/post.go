package user

import (
	keyUtil "dgram/modules/util"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Text     *string   `json:"text"`
	Images   *[]string `json:"images"`
	Videos   *[]string `json:"videos"`
	Comments *[]string `json:"comments"`
	Votes    *[]Vote   `json:"votes"`
	gorm.Model
}

type Comment struct {
	ID    uuid.UUID `json:"id"`
	Text  string    `json:"text"`
	Votes []Vote    `json:"votes"`
}

type Vote struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func CreateNewPost(c *fiber.Ctx) (string, error) {

	UploadPostMedia(c)

	var NewPost Post
	user := FindUser(c.Params("id"))

	NewPost.ID = uuid.New()
	NewPost.UserID = user.ID

	//user.Posts = append(user.Posts,)
	error := c.BodyParser(&NewPost)
	if error != nil || isValidPost(NewPost) {
		return "", FailedTransaction(c)
	}

	return "TANGLEADDRESS", nil
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

func isValidPost(p Post) bool {

	if len(*p.Images) == 0 && len(*p.Videos) == 0 {
		return false
	}

	return true
}
