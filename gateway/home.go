package gateway

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

type addHomeReq struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// AddHome create and add user to an home
func AddHome(c echo.Context) error {
	req := new(addHomeReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	var missingFields []string
	if req.Name == "" {
		missingFields = append(missingFields, "name")
	}
	if req.Address == "" {
		missingFields = append(missingFields, "address")
	}
	if len(missingFields) > 0 {
		return c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "Some fields missing: " + strings.Join(missingFields, ", "),
		})
	}

	// var user User
	// err := DB.Get(&user, "SELECT * FROM user WHERE email=$1", req.Email)
	// if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
	// 	return c.JSON(http.StatusBadRequest, MessageResponse{
	// 		Message: "Email and password doesn't match",
	// 	})
	// }

	// id := NewULID().String()
	// createdAt := time.Now()

	// newToken := Token{
	// 	ID:        id,
	// 	UserID:    user.ID,
	// 	Type:      "signin",
	// 	IP:        c.RealIP(),
	// 	UserAgent: c.Request().UserAgent(),
	// 	Read:      1,
	// 	Write:     1,
	// 	Manage:    1,
	// 	Admin:     1,
	// 	CreatedAt: createdAt.String(),
	// 	ExpireAt:  createdAt.AddDate(0, 1, 0).String(),
	// }
	// DB.NamedExec("INSERT INTO token (id, user_id, type, ip, user_agent, read, write, manage, admin, created_at, expire_at) VALUES (:id, :user_id, :type, :ip, :user_agent, :read, :write, :manage, :admin, :created_at, :expire_at)", newToken)

	return c.JSON(http.StatusOK, MessageResponse{
		Message: "test",
	})
}
