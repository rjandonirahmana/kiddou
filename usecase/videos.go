package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kiddou/domain"
	"log"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

type usecaseVideos struct {
	repoVides domain.RepoVideos
	db        *sql.DB
}

func NewUsecaseVideos(repoVideos domain.RepoVideos, db *sql.DB) *usecaseVideos {
	return &usecaseVideos{repoVides: repoVideos, db: db}
}

const allowedExtVideos = ".mp4"

func (u *usecaseVideos) InsertVideos(ctx context.Context, videos []byte, input *domain.InserVideosRequest) error {
	mime := mimetype.Detect(videos)
	if !strings.Contains(allowedExtVideos, mime.Extension()) {
		return errors.New("File Type is not allowed, file type: " + mime.Extension())
	}

	log.Println("sampe sini 11")
	availableCategories, err := u.repoVides.GetAllCategories(ctx)
	if err != nil {
		return err
	}
	found := false
	for _, v := range availableCategories {
		if v.ID == input.CategoryID {
			found = true
		}
	}
	if !found {
		return errors.New("category not found")
	}
	path := "/home/rjandoni/Desktop/kiddou/videos/"

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	vid := &domain.Videos{
		CategoryID:   input.CategoryID,
		Name:         input.Name,
		Descriptions: input.Descriptions,
		Price:        input.Price,
		Url:          "/home/rjandoni/Desktop/kiddou/videos/",
		Subscribers:  0,
	}
	folder, err := os.Create(fmt.Sprintf("%s%d%s", path, 1, mime.Extension()))
	if err != nil {
		return err
	}
	defer folder.Close()

	_, err = folder.Write(videos)
	if err != nil {
		return err
	}

	tx, err := u.db.Begin()
	if err != nil {
		return err
	}

	err = u.repoVides.InsertVideos(ctx, tx, vid)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil

}

func (u *usecaseVideos) SubscribtionVideo(ctx context.Context, userID string, videoID int) error {

	return nil
}
