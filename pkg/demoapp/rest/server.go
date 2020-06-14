package restv1

import (
	"net/http"

	"github.com/hippeus/appbase/pkg/buildinfo"
	"github.com/hippeus/appbase/pkg/logger"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	logger.Logger
	AppName string
	buildinfo.Build
}

// LivenessProbeRequest GET /api/v1/demoapp/.well-known/alive
func (h *Handler) LivenessProbeRequest(ctx echo.Context, params LivenessProbeRequestParams) error {
	res := LivenessProbeResponse{}
	if params.Full != nil {
		res = h.Build.VersionAsMap()
		res["application"] = h.AppName
	}

	if len(res) == 0 {
		return ctx.NoContent(http.StatusOK)
	}

	return ctx.JSON(http.StatusOK, &res)
}
