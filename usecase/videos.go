package usecase

import "kiddou/domain"

type usecaseVideos struct {
	repoVides domain.RepoVideos
}

func NewUsecaseVideos(repoVideos domain.RepoVideos) *usecaseVideos {
	return &usecaseVideos{repoVides: repoVideos}
}
