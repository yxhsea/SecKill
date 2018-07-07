package product

import (
	"SecKill/sk_admin/model"
	"SecKill/sk_admin/service"
	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
	"log"
)

func CreateProduct(ctx *gin.Context) {
	var product = &model.Product{}
	product.ProductName = ctx.PostForm("product_name")
	product.Total, _ = com.StrTo(ctx.PostForm("product_total")).Int()
	product.Status, _ = com.StrTo(ctx.PostForm("status")).Int()

	productServer := service.NewProductServer()
	err := productServer.CreateProduct(product)
	if err != nil {
		log.Printf("ProductServer.CreateProduct, err : %v", err)
		ctx.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "failed",
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "success",
	})
	return
}

func GetPorductList(ctx *gin.Context) {
	productService := service.NewProductServer()
	productList, err := productService.GetProductList()
	if err != nil {
		log.Printf("ProductService.productList, err : %v", err)
		ctx.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "failed",
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": productList,
	})
	return
}
