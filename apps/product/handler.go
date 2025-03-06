package product

import (
	"Ecommerce-basic/infra/gin"
	"Ecommerce-basic/infra/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	svc service
}

func newHandler(svc service) handler {
	return handler{
		svc: svc,
	}
}

func (h handler) CreateProduct(ctx *gin.Context) {
	var req CreateProductRequestPayload

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := infragin.NewResponse(
			infragin.WithHttpCode(http.StatusBadRequest),
			infragin.WithMessage("invalid payload"),
			infragin.WithError(response.ErrorBadRequest),
		)
		resp.Send(ctx)
		return
	}

	if err := h.svc.CreateProduct(ctx, req); err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		resp := infragin.NewResponse(
			infragin.WithMessage(err.Error()),
			infragin.WithError(myErr),
		)
		resp.Send(ctx)
		return
	}

	resp := infragin.NewResponse(
		infragin.WithHttpCode(http.StatusCreated),
		infragin.WithMessage("create product success"),
	)
	resp.Send(ctx)
}

func (h handler) GetListProducts(ctx *gin.Context) {
	var req ListProductRequestPayload

	if err := ctx.ShouldBindQuery(&req); err != nil {
		resp := infragin.NewResponse(
			infragin.WithHttpCode(http.StatusBadRequest),
			infragin.WithMessage("invalid payload"),
			infragin.WithError(response.ErrorBadRequest),
		)
		resp.Send(ctx)
		return
	}

	products, err := h.svc.ListProducts(ctx, req)
	if err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		resp := infragin.NewResponse(
			infragin.WithMessage(err.Error()),
			infragin.WithError(myErr),
		)
		resp.Send(ctx)
		return
	}

	productListResponse := NewProductListResponseFromEntity(products)

	resp := infragin.NewResponse(
		infragin.WithHttpCode(http.StatusOK),
		infragin.WithMessage("get list products success"),
		infragin.WithPayload(productListResponse),
		infragin.WithQuery(req.GenerateDefaultValue()),
	)
	resp.Send(ctx)
}

func (h handler) GetProductDetail(ctx *gin.Context) {
	sku := ctx.Param("sku")
	if sku == "" {
		resp := infragin.NewResponse(
			infragin.WithHttpCode(http.StatusBadRequest),
			infragin.WithMessage("invalid SKU"),
			infragin.WithError(response.ErrorBadRequest),
		)
		resp.Send(ctx)
		return
	}

	product, err := h.svc.ProductDetail(ctx, sku)
	if err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		resp := infragin.NewResponse(
			infragin.WithMessage(err.Error()),
			infragin.WithError(myErr),
		)
		resp.Send(ctx)
		return
	}

	productDetail := ProductDetailResponse{
		Id:        product.Id,
		Name:      product.Name,
		SKU:       product.SKU,
		Stock:     product.Stock,
		Price:     product.Price,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}

	resp := infragin.NewResponse(
		infragin.WithHttpCode(http.StatusOK),
		infragin.WithMessage("get product detail success"),
		infragin.WithPayload(productDetail),
	)
	resp.Send(ctx)
}
