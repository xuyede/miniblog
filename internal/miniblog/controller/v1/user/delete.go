package user

import (
	"github.com/gin-gonic/gin"
	"github.com/xuyede/miniblog/internal/pkg/core"
	"github.com/xuyede/miniblog/internal/pkg/log"
)

func (ctrl *UserController) Delete(c *gin.Context) {
	log.C(c).Infow("DELETE /v1/users/:name called")

	username := c.Param("name")
	if err := ctrl.b.Users().Delete(c, username); err != nil {
		core.GenarateResponse(c, err, nil)
		return
	}

	if _, err := ctrl.a.RemoveNamedPolicy("p", username, "", ""); err != nil {
		core.GenarateResponse(c, err, nil)
		return
	}

	core.GenarateResponse(c, nil, nil)
}
