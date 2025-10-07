package factory

import "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"

type ItemFactory interface {
	CreateItem() models.ItemInterface
}

type ItemFactoryConcrete struct{}

func NewItemFactoryConcrete() ItemFactory {
	return &ItemFactoryConcrete{}
}

// CreateItem implements ItemFactory.
func (i *ItemFactoryConcrete) CreateItem() models.ItemInterface {
	panic("unimplemented")
}
