package usecase

import (
	"context"

	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
)

func (s Service) Upsert(ctx context.Context, req model.SiteRes) error {
	return s.repo.UpSertFranchiseSite(ctx, req)
}
