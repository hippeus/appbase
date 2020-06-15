package spa

import (
	"bytes"
	"io"
	"net/http"

	"github.com/elliots/litter"
	"github.com/hippeus/appbase/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/markbates/pkger"
)

//go:generate pkger

const (
	// DistPath describe targeted root directory for embedded files. Should
	// treat GOMOD as filesystem root.
	// Must be static and known at compile time so pkger tool could parse it.
	DistPath = "/ui/build"
)

type SPA struct {
	MountPath string

	logger.Logger
	embedFSRoot pkger.Dir
}

func (s *SPA) MountToEchoRouter(router *echo.Echo, m ...echo.MiddlewareFunc) {
	if s.embedFSRoot == "" {
		s.embedFSRoot = pkger.Dir(DistPath)
	}
	if s.Logger == nil {
		s.Logger = logger.NOP{}
	}

	if s.MountPath == "" {
		s.MountPath = "/*"
	}

	info, _ := pkger.Current()
	litter.DumpColor(info)
	s.Infof("Mount SPA to %s", s.MountPath)
	_ = router.GET(s.MountPath, s.echoHandleEmbeddedFiles, m...)
}

func (s *SPA) echoHandleEmbeddedFiles(ctx echo.Context) error {
	const (
		indexPage = "index.html"
	)

	uri := ctx.Request().URL.Path
	if uri == "/" {
		uri = indexPage
	}

	file, err := s.embedFSRoot.Open(uri)
	if err != nil {
		// fallback to the landing page
		file, err = s.embedFSRoot.Open(indexPage)
		if err != nil {
			s.Errorf("rolling back to the landing page failed")
			return ctx.NoContent(http.StatusNotFound)
		}
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, file)
	if err != nil {
		info, _ := file.Stat()
		s.Errorf("failed to load file %s, for url request: %s", info.Name, uri)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load a file").SetInternal(err)
	}

	ct := http.DetectContentType(buf.Bytes())
	return ctx.Blob(http.StatusOK, ct, buf.Bytes())
}
