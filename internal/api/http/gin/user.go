package gin

import (
	"net/http"
	"todo-app/domain"
	"todo-app/pkg/clients"
	"todo-app/pkg/tokenprovider"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	Register(data *domain.UserCreate) error
	Login(data *domain.UserLogin) (tokenprovider.Token, error)
	GetByIdUser(id uuid.UUID) (domain.User, error)
	UpdateUser(id uuid.UUID, data *domain.UserUpdate) error
	DeleteUser(id uuid.UUID) error
}

type userHandler struct {
	userService UserService
}

func NewUserHandler(apiVersion *gin.RouterGroup, svc UserService, middlewareAuth gin.HandlerFunc) {
	userHandler := &userHandler{
		userService: svc,
	}

	users := apiVersion.Group("/users")
	users.POST("/register", userHandler.RegisterUserHandler)
	users.POST("/login", userHandler.LoginHandler)
	users.Use(middlewareAuth)
	users.GET("/:id", userHandler.GetUserhandler)
	users.PATCH("/:id", userHandler.UpdateUserHandler)
	users.DELETE("/:id", userHandler.DeleteUserHandler)
}

// RegisterUserHandler godoc
// @Summary Đăng ký người dùng mới
// @Description Tạo người dùng mới trong hệ thống
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.UserCreate true "Thông tin người dùng"
// @Success 201 {object} domain.User
// @Router /users/register [post]
func (h *userHandler) RegisterUserHandler(c *gin.Context) {
	var data domain.UserCreate

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))
		return
	}

	if err := h.userService.Register(&data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, clients.SimpleSuccessResponse(data.ID))
}

// LoginHandler godoc
// @Summary Đăng nhập
// @Description Đăng nhập người dùng và nhận token
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.UserLogin true "Thông tin đăng nhập"
// @Success 200 {object} tokenprovider.Token
// @Router /users/login [post]
func (h *userHandler) LoginHandler(c *gin.Context) {
	var data domain.UserLogin

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))
		return
	}

	token, err := h.userService.Login(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(token))
}

// GetUserhandler godoc
// @Summary Lấy thông tin người dùng theo ID
// @Description Lấy thông tin người dùng dựa trên ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID người dùng" // Chỉnh sửa ở đây
// @Success 200 {object} domain.User
// @Router /users/{id} [get]
func (h *userHandler) GetUserhandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))
		return
	}

	user, err := h.userService.GetByIdUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUserHandler godoc
// @Summary Cập nhật thông tin người dùng
// @Description Cập nhật thông tin người dùng theo ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID người dùng" // Chỉnh sửa ở đây
// @Param user body domain.UserUpdate true "Thông tin người dùng"
// @Success 200 {object} domain.User
// @Router /users/{id} [patch]
func (h *userHandler) UpdateUserHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))
		return
	}

	data := domain.UserUpdate{}
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))
		return
	}

	if err := h.userService.UpdateUser(id, &data); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(id))
}

// DeleteUserHandler godoc
// @Summary Xóa người dùng theo ID
// @Description Xóa người dùng khỏi hệ thống
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID người dùng" // Chỉnh sửa ở đây
// @Router /users/{id} [delete]
func (h *userHandler) DeleteUserHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusNoContent, nil) // Sử dụng http.StatusNoContent cho phản hồi thành công khi xóa
}
