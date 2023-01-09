package v1

import (
	"aggregate-task/internal/busi/core"
	"aggregate-task/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Walk tipsets godoc
// @Description walk the historical DAG's tipsets.
// @Tags Aggregate-Task-API-Internal-V1-CallByManual
// @Accept application/json,json
// @Produce application/json,json
// @Param Walk query core.Walk false "Walk"
// @Success 200 {object} nil
// @Failure 400 {object} utils.ResponseWithRequestId
// @Failure 500 {object} utils.ResponseWithRequestId
// @Router /api/v1/walk [post]
func WalkTipsets(c *gin.Context) {
	app := utils.Gin{C: c}

	var r core.Walk
	if err := c.ShouldBindQuery(&r); err != nil {
		app.HTTPResponse(http.StatusBadRequest, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	if err := r.Validate(); err != nil {
		app.HTTPResponse(http.StatusBadRequest, utils.NewResponse(utils.CodeBadRequest, err.Error(), nil))
		return
	}

	task, _ := c.Get(TASK)
	r.Task, _ = task.(string)

	dependentTasks, _ := c.Get(DEPENDENTTASKS)
	r.DependentTasks, _ = dependentTasks.([]string)

	resp := core.WalkTipsetsRun(c.Request.Context(), &r)
	if resp != nil {
		app.HTTPResponse(resp.HttpCode, resp.Response)
		return
	}

	app.HTTPResponseOK(nil)
}
