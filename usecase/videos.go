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
	"time"

	"github.com/gabriel-vasile/mimetype"
)

type usecaseVideos struct {
	repoVides    domain.RepoVideos
	subcribeRepo domain.SubscribersRepo
	db           *sql.DB
}

func NewUsecaseVideos(repoVideos domain.RepoVideos, db *sql.DB, repoSubscribe domain.SubscribersRepo) *usecaseVideos {
	return &usecaseVideos{repoVides: repoVideos, db: db, subcribeRepo: repoSubscribe}
}

const allowedExtVideos = ".mp4"

func (u *usecaseVideos) InsertVideos(ctx context.Context, videos []byte, input *domain.InserVideosRequest) error {
	mime := mimetype.Detect(videos)
	if !strings.Contains(allowedExtVideos, mime.Extension()) {
		return errors.New("File Type is not allowed, file type: " + mime.Extension())
	}

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

	tx, err := u.db.Begin()
	if err != nil {
		return err
	}

	video, err := u.repoVides.GetByID(ctx, videoID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if video.ID == 0 {
		return errors.New("video not found")
	}

	sub, err := u.subcribeRepo.GetByVideoID(ctx, videoID, userID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	log.Println(sub)
	if sub.ID != 0 {
		if sub.Status == "expired" {
			video.Subscribers += 1
			err = u.repoVides.UpdateVideo(ctx, tx, &video)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		sub.ExpiredAt = time.Now().AddDate(0, 1, 0)
		sub.Status = "active"
		err = u.subcribeRepo.UpdateSub(ctx, tx, &sub)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {

		video.Subscribers += 1
		err = u.repoVides.UpdateVideo(ctx, tx, &video)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = u.subcribeRepo.InsertSub(ctx, tx, &domain.Subscribers{UserID: userID, VideoID: videoID, TypeSubscription: "freemium", ExpiredAt: time.Now().AddDate(0, 1, 0), Status: "active", SubscribeAT: time.Now()})
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (u *usecaseVideos) SubsribesStatus(ctx context.Context, userID string, videoID int) (res domain.SubsriptionStatusDTO, err error) {
	sub, err := u.subcribeRepo.GetByVideoID(ctx, videoID, userID)
	if err != nil && sql.ErrNoRows != nil {
		return
	}

	if sub.ID == 0 {
		return res, errors.New("you havent subscribe this video")
	}

	video, err := u.repoVides.GetByID(ctx, videoID)
	if err != nil {
		return
	}

	res.VideoName = video.Name
	res.TypeSubscription = sub.TypeSubscription
	res.ExpiredAt = sub.ExpiredAt.Format(time.RFC3339)
	res.Status = sub.Status
	res.SubscribeAT = sub.SubscribeAT.Format(time.RFC3339)

	return

}

func (u *usecaseVideos) GetByCategory(ctx context.Context, categoryID int) (res []domain.Videos, err error) {
	availableCategories, err := u.repoVides.GetAllCategories(ctx)
	if err != nil {
		return res, err
	}
	found := false
	for _, v := range availableCategories {
		if v.ID == categoryID {
			found = true
		}
	}
	if !found {
		return res, errors.New("category not found")
	}

	res, err = u.repoVides.GetByCategory(ctx, categoryID)
	if err != nil {
		return res, err
	}
	return

}
