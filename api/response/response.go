package response

import "go.uber.org/zap"

type Responser struct {
	logger *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger) Responser {
	return Responser{
		logger,
	}
}
