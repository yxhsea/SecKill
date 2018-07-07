package service

import (
	"SecKill/sk_admin/model"
	"log"
)

type ProductServer struct {
}

func NewProductServer() *ProductServer {
	return &ProductServer{}
}

func (p *ProductServer) CreateProduct(product *model.Product) error {
	productEntity := model.NewProductModel()
	err := productEntity.CreateProduct(product)
	if err != nil {
		log.Printf("ProductEntity.CreateProduct, err : %v", err)
		return err
	}
	return nil
}

func (p *ProductServer) GetProductList() ([]map[string]interface{}, error) {
	productEntity := model.NewProductModel()
	productList, err := productEntity.GetProductList()
	if err != nil {
		log.Printf("ProductEntity.CreateProduct, err : %v", err)
		return nil, err
	}
	return productList, nil
}
