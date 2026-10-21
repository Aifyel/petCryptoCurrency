package gin

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Aifyel/petCryptoCurrency/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Aifyel/petCryptoCurrency/docs"
)

type StatisticsService interface {
	GetStatistics(ctx context.Context, currencies []string) ([]entities.CurrencyStatistics, error)
}

type Server struct {
	serv    *http.Server
	service StatisticsService
	logger  *zap.Logger
}

func NewServer(addr string, timeout time.Duration, service StatisticsService, logger *zap.Logger) (*Server, error) {
	if addr == "" {
		return nil, errors.Wrap(entities.ErrInvalidParams, "gin: empty server address")
	}
	if timeout <= 0 {
		return nil, errors.Wrap(entities.ErrInvalidParams, "gin: invalid timeout")
	}
	if service == nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "gin: service cannot be nil")
	}
	if logger == nil {
		return nil, errors.Wrap(entities.ErrInvalidParams, "gin: logger cannot be nil")
	}

	s := &Server{
		serv: &http.Server{
			Addr:              addr,
			ReadHeaderTimeout: timeout,
		},
		logger:  logger,
		service: service,
	}

	s.setupRouter()

	return s, nil
}

func (s *Server) setupRouter() {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/stats", s.GetStats)

	s.serv.Handler = router
}

func (s *Server) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	go s.Start(errChan)

	select {
	case <-ctx.Done():
		s.logger.Info("shutting down http server gracefully")

		ctxSD, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := s.Shutdown(ctxSD); err != nil {
			return errors.Wrap(err, "http server shutdown failed")
		}
		s.logger.Info("http server successfully stopped")
	case err := <-errChan:
		return errors.Wrap(err, "http server failed")
	}
	return nil
}

func (s *Server) Start(errChan chan<- error) {
	s.logger.Info("starting http server", zap.String("address", s.serv.Addr))
	if err := s.serv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errChan <- err
	}
	close(errChan)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.serv.Shutdown(ctx)
}

// GetStats godoc
// @Summary Получить статистику криптовалют
// @Description Возвращает текущую, минимальную, максимальную цену и процент изменения за час для указанных криптовалют.
// @Tags stats
// @Accept json
// @Produce json
// @Param currencies query string true "Список криптовалют через запятую (например: bitcoin,ethereum)"
// @Success 200 {array} entities.CurrencyStatistics
// @Failure 400 {object} map[string]string "Ошибка валидации параметров"
// @Failure 502 {object} map[string]string "Внешний сервис недоступен"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /stats [get]
func (s *Server) GetStats(c *gin.Context) {
	raw := c.Query("currencies")

	currencies := strings.Split(raw, ",")
	for i := 0; i < len(currencies); i++ {
		currencies[i] = strings.ToLower(strings.TrimSpace(currencies[i]))
		if currencies[i] == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "currencies must be a comma-separated list of non-empty values"})
			return
		}
	}

	ctx := c.Request.Context()
	stats, err := s.service.GetStatistics(ctx, currencies)
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrInvalidParams):
			s.logger.Warn("invalid request", zap.Error(err), zap.Strings("currencies", currencies))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request parameters"})
		case errors.Is(err, entities.ErrClientFailure):
			s.logger.Error("upstream failure", zap.Error(err), zap.Strings("currencies", currencies))
			c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service unavailable"})
		default:
			s.logger.Error("failed to get statistics", zap.Error(err), zap.Strings("currencies", currencies))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.JSON(http.StatusOK, stats)
}
