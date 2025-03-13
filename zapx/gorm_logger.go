package zapx

import (
	"context"
	"fmt"
	"gorm.io/gorm/logger"
	"time"
)

func NewGormLoggerAdapter(logger LoggerX) *GormLoggerAdapter {
	return &GormLoggerAdapter{logger: logger}
}

type GormLoggerAdapter struct {
	logger LoggerX
}

func (g *GormLoggerAdapter) LogMode(level logger.LogLevel) logger.Interface {
	//TODO implement me
	panic("implement me")
}

func (g *GormLoggerAdapter) Info(ctx context.Context, s string, i ...interface{}) {
	g.logger.WithContext(ctx).Value("gorm")
	g.logger.Infof(s, fmt.Sprintf("%v", i))
}

func (g *GormLoggerAdapter) Warn(ctx context.Context, s string, i ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (g *GormLoggerAdapter) Error(ctx context.Context, s string, i ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (g *GormLoggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	//TODO implement me
	panic("implement me")
}
