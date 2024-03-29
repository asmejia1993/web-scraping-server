package repository

import (
	"context"

	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
)

type IFranchiseRepository interface {
	FindFranchisesById(id string, ctx context.Context) model.FranchiseInfo
	CreateFranchisesHotel(ctx context.Context, req model.FranchiseInfo) (string, error)
	UpSertFranchiseSite(ctx context.Context, site model.SiteRes) error
	All(ctx context.Context, params map[string][]string) ([]model.FranchiseInfo, error)
}
