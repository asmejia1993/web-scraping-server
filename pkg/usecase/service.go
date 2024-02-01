package usecase

import "github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/repository"

type Service struct {
	repo      repository.IFranchiseRepository
	redisRepo repository.IRedisRepository
}

func NewService(r repository.IFranchiseRepository, redis repository.IRedisRepository) Service {
	return Service{
		repo:      r,
		redisRepo: redis,
	}
}
