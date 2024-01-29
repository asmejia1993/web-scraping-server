package usecase

import "github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/repository"

type Service struct {
	repo repository.IFranchiseRepository
}

func NewService(r repository.IFranchiseRepository) Service {
	return Service{
		repo: r,
	}
}
