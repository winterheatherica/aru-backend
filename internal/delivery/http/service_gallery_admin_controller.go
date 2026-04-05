package http

import (
	"mime/multipart"
	"strconv"
	"strings"

	"aru-backend/internal/model"
	adminusecase "aru-backend/internal/usecase/admin"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ServiceGalleryAdminController struct {
	Usecase adminusecase.ServiceGalleryUsecase
}

func NewServiceGalleryAdminController(u adminusecase.ServiceGalleryUsecase) *ServiceGalleryAdminController {
	return &ServiceGalleryAdminController{Usecase: u}
}

func (c *ServiceGalleryAdminController) ListByService(ctx *fiber.Ctx) error {
	service := strings.TrimSpace(ctx.Params("service"))
	lang := strings.ToUpper(strings.TrimSpace(ctx.Query("lang", "ID")))
	items, err := c.Usecase.ListByService(ctx.Context(), service, lang)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(items)
}

func (c *ServiceGalleryAdminController) GetByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	lang := strings.ToUpper(strings.TrimSpace(ctx.Query("lang", "ID")))
	item, err := c.Usecase.GetByID(ctx.Context(), id, lang)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.JSON(item)
}

func (c *ServiceGalleryAdminController) Create(ctx *fiber.Ctx) error {
	input, image, err := parseServiceGalleryUpsertInput(ctx, true)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	created, err := c.Usecase.Create(ctx.Context(), input, image)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.Status(fiber.StatusCreated).JSON(created)
}

func (c *ServiceGalleryAdminController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	input, image, err := parseServiceGalleryUpsertInput(ctx, false)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	updated, err := c.Usecase.Update(ctx.Context(), id, input, image)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.JSON(updated)
}

func (c *ServiceGalleryAdminController) HardDelete(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	if err := c.Usecase.HardDelete(ctx.Context(), id); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.JSON(fiber.Map{"message": "deleted"})
}

func parseServiceGalleryUpsertInput(ctx *fiber.Ctx, imageRequired bool) (model.ServiceGalleryUpsertInput, *multipart.FileHeader, error) {
	activeRaw := strings.TrimSpace(ctx.FormValue("is_active", "true"))
	active := true
	if activeRaw != "" {
		v, perr := strconv.ParseBool(activeRaw)
		if perr != nil {
			return model.ServiceGalleryUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "is_active must be boolean")
		}
		active = v
	}

	input := model.ServiceGalleryUpsertInput{
		Service:   strings.ToUpper(strings.TrimSpace(ctx.FormValue("service"))),
		MediaType: strings.ToUpper(strings.TrimSpace(ctx.FormValue("media_type", "IMAGE"))),
		Language:  strings.ToUpper(strings.TrimSpace(ctx.FormValue("language", "ID"))),
		Title:     ptrOrNil(ctx.FormValue("title")),
		Alt:       ptrOrNil(ctx.FormValue("alt")),
		Caption:   ptrOrNil(ctx.FormValue("caption")),
		IsActive:  active,
	}

	image, _ := ctx.FormFile("image")
	if imageRequired && image == nil {
		return model.ServiceGalleryUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "image is required")
	}
	return input, image, nil
}
