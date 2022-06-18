package handler

import (
	"context"
	"errors"
	"kiddou/domain"
	"kiddou/grpc/videos"
	"log"
	"sync"
)

type GrcpServicesVideos struct {
	videos.UnimplementedVideosStreanServer
	usecaseVideo domain.UsecaseVideos
}

func NewGrpcVideos(usecaseVideo domain.UsecaseVideos) *GrcpServicesVideos {
	return &GrcpServicesVideos{usecaseVideo: usecaseVideo}
}

func (g *GrcpServicesVideos) GetByCategoryID(req *videos.VideosRequest, srv videos.VideosStrean_GetByCategoryIDServer) error {

	log.Printf("fetch response for id : %d", req.CategoryID)
	ctx := context.TODO()

	//use wait group to allow process to be concurrent
	vids, err := g.usecaseVideo.GetByCategory(ctx, int(req.CategoryID))
	if err != nil {
		return err
	}
	if len(vids) == 0 {
		return errors.New("data video not found")
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for _, v := range vids {

			//time sleep to simulate server process time

			resp := videos.VideoResponse{Id: int32(v.ID), CategoryID: int32(v.CategoryID), Name: v.Name, Description: v.Descriptions, Price: v.Price, Url: v.Url, Subscriber: int32(v.Subscribers)}
			if err := srv.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}
			log.Printf("finishing request number")

		}
		defer wg.Done()
	}()

	wg.Wait()
	return nil
}
