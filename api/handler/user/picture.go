package user

import (
	"io"
	"net/http"

	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/api/request"
	"github.com/werdna521/userland/api/response"
	"github.com/werdna521/userland/service"
)

const fileLimit = 200 * 1024

type setProfilePictureResponse struct {
	Success bool `json:"success"`
}

func SetProfilePicture(us service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > fileLimit {
			response.Error(w, e.NewRequestEntityTooLargeError("file too large")).JSON()
			return
		}

		r.ParseMultipartForm(200 * 1024 * 1024)
		file, _, err := r.FormFile("file")
		if err != nil {
			response.Error(w, e.NewBadRequestError("cannot parse file")).JSON()
			return
		}
		defer file.Close()

		var fileHeader = make([]byte, 512)
		if _, err := file.Read(fileHeader); err != nil {
			response.Error(w, e.NewBadRequestError("cannot read file")).JSON()
			return
		}

		// set position back to start.
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			response.Error(w, e.NewInternalServerError())
			return
		}

		if http.DetectContentType(fileHeader) != "image/png" {
			response.Error(w, e.NewBadRequestError("picture must be a .png file")).JSON()
			return
		}

		ctx := r.Context()
		at, err := request.GetAccessTokenFromCtx(ctx)
		if err != nil {
			response.Error(w, err.(e.Error)).JSON()
			return
		}

		err = us.SetProfilePicture(ctx, at.UserID, file)
		if err != nil {
			response.Error(w, err.(e.Error)).JSON()
			return
		}

		response.OK(w, &setProfilePictureResponse{
			Success: true,
		}).JSON()
	}
}
