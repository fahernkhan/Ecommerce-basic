package product

import (
	"Ecommerce-basic/infra/gin"
	"Ecommerce-basic/infra/response"
	"net/http"
	"strconv"

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

func (h handler) UpdateProduct(c *gin.Context) {
	var req UpdateProductRequestPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := infragin.NewResponse(
			infragin.WithHttpCode(http.StatusBadRequest),
			infragin.WithMessage("invalid payload"),
			infragin.WithError(response.ErrorBadRequest),
		)
		resp.Send(c)
		return
	}

	id := c.Param("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		resp := infragin.NewResponse(
			infragin.WithHttpCode(http.StatusBadRequest),
			infragin.WithMessage("invalid product ID"),
			infragin.WithError(response.ErrorBadRequest),
		)
		resp.Send(c)
		return
	}

	if err := h.svc.UpdateProduct(c.Request.Context(), productID, req); err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		resp := infragin.NewResponse(
			infragin.WithMessage(err.Error()),
			infragin.WithError(myErr),
		)
		resp.Send(c)
		return
	}

	resp := infragin.NewResponse(
		infragin.WithHttpCode(http.StatusOK),
		infragin.WithMessage("update product success"),
	)
	resp.Send(c)
}

func (h handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		resp := infragin.NewResponse(
			infragin.WithHttpCode(http.StatusBadRequest),
			infragin.WithMessage("invalid product ID"),
			infragin.WithError(response.ErrorBadRequest),
		)
		resp.Send(c)
		return
	}

	if err := h.svc.DeleteProduct(c.Request.Context(), productID); err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		resp := infragin.NewResponse(
			infragin.WithMessage(err.Error()),
			infragin.WithError(myErr),
		)
		resp.Send(c)
		return
	}

	resp := infragin.NewResponse(
		infragin.WithHttpCode(http.StatusOK),
		infragin.WithMessage("delete product success"),
	)
	resp.Send(c)
}

func (h handler) SearchProducts(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		resp := infragin.NewResponse(
			infragin.WithHttpCode(http.StatusBadRequest),
			infragin.WithMessage("keyword is required"),
			infragin.WithError(response.ErrorBadRequest),
		)
		resp.Send(c)
		return
	}

	pagination := ProductPagination{
		Cursor: 0,
		Size:   10,
	}

	products, err := h.svc.SearchProducts(c.Request.Context(), keyword, pagination)
	if err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		resp := infragin.NewResponse(
			infragin.WithMessage(err.Error()),
			infragin.WithError(myErr),
		)
		resp.Send(c)
		return
	}

	productListResponse := NewProductListResponseFromEntity(products)

	resp := infragin.NewResponse(
		infragin.WithHttpCode(http.StatusOK),
		infragin.WithMessage("search products success"),
		infragin.WithPayload(productListResponse),
	)
	resp.Send(c)
}

func (h handler) FilterProducts(c *gin.Context) {
	minPrice, _ := strconv.Atoi(c.Query("minPrice"))
	maxPrice, _ := strconv.Atoi(c.Query("maxPrice"))
	minStock, _ := strconv.Atoi(c.Query("minStock"))
	maxStock, _ := strconv.Atoi(c.Query("maxStock"))

	pagination := ProductPagination{
		Cursor: 0,
		Size:   10,
	}

	products, err := h.svc.FilterProducts(c.Request.Context(), minPrice, maxPrice, int16(minStock), int16(maxStock), pagination)
	if err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}
		resp := infragin.NewResponse(
			infragin.WithMessage(err.Error()),
			infragin.WithError(myErr),
		)
		resp.Send(c)
		return
	}

	productListResponse := NewProductListResponseFromEntity(products)

	resp := infragin.NewResponse(
		infragin.WithHttpCode(http.StatusOK),
		infragin.WithMessage("filter products success"),
		infragin.WithPayload(productListResponse),
	)
	resp.Send(c)
}
