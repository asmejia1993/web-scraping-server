package usecase

import "context"

func (s Service) GetKey(ctx context.Context, key string) (string, error) {
	return s.redisRepo.Get(ctx, key)
}

func (s Service) SetKey(ctx context.Context, key string, value interface{}) {
	s.redisRepo.Set(ctx, key, value)
}

func (s Service) Exist(ctx context.Context, key string) (int64, error) {
	return s.redisRepo.Exist(ctx, key)
}
