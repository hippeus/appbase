package spa

import (
	"net/http"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/hippeus/appbase/pkg/logger"
	"github.com/labstack/echo/v4"
)

type SPA struct {
	logger.Logger

	Box    *packr.Box
	Prefix string
}

var (
	FrontendApp = &SPA{
		Box: packr.New("demoapp", "../../ui/build"),
	}
)

func (s *SPA) MountToEchoRouter(router *echo.Echo, m ...echo.MiddlewareFunc) {
	if s.Prefix == "" {
		s.Prefix = "/*"
	}
	_ = router.GET(s.Prefix, s.echoHandleEmbeddedFiles, m...)
}

func (s *SPA) echoHandleEmbeddedFiles(ctx echo.Context) error {
	const (
		indexPage = "index.html"
	)

	uri := ctx.Request().URL.Path
	if uri == "/" {
		uri = indexPage
	}

	uri = strings.TrimPrefix(uri, s.Prefix)

	body, err := s.Box.Find(uri)
	if err != nil {
		// fallback to the landing page
		body, err = s.Box.Find(indexPage)
		if err != nil {
			s.Errorf("rolling back to the landing page failed")
			return ctx.NoContent(http.StatusNotFound)
		}
	}

	ct := http.DetectContentType(body)
	return ctx.Blob(http.StatusOK, ct, body)
}
