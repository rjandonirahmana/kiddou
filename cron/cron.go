package cron

import (
	"context"
	"database/sql"
	"kiddou/domain"
	"time"
	// "github.com/jasonlvhit/gocron"
)

// func Start(db *sql.DB, repoSub domain.SubscribersRepo) error {
// 	s := gocron.NewScheduler()
// 	gocron.Clear()

// 	err := s.Every(30).Seconds().Do(TaskMonitoring(db, repoSub))
// 	if err != nil {
// 		return err

// 	}
// 	s.Start()
// 	return nil
// }

func TaskMonitoring(db *sql.DB, repoSub domain.SubscribersRepo, repoVideo domain.RepoVideos) error {
	querry := `SELECT id, user_id, video_id, type_subscription, expired_at, status FROM subcriptions WHERE status = $1`
	ctx := context.Background()

	row, err := db.QueryContext(ctx, querry, "active")
	if err != nil {
		return err
	}

	var subs []domain.Subscribers
	defer row.Close()
	for row.Next() {
		sub := domain.Subscribers{}
		err = row.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.VideoID,
			&sub.TypeSubscription,
			&sub.ExpiredAt,
			&sub.Status,
		)
		if err != nil {
			return err
		}

		subs = append(subs, sub)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, s := range subs {
		if time.Now().After(s.ExpiredAt) {
			s.Status = "expired"
			err = repoSub.UpdateSub(ctx, tx, &s)
			if err != nil {
				tx.Rollback()
				return err
			}

			video, err := repoVideo.GetByID(ctx, s.VideoID)
			if err != nil {
				tx.Rollback()
				return err
			}
			video.Subscribers -= 1
			err = repoVideo.UpdateVideo(ctx, tx, &video)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	return nil
}
