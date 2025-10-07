package factory

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	ItemSetRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_set"
)

type ItemFactoryWithChildrenConcrete struct {
	itemRepo    ItemRepo.Repository
	itemSetRepo ItemSetRepo.Repository
}

func NewItemFactoryWithChildrenConcrete(itemRepo ItemRepo.Repository, itemSetRepo ItemSetRepo.Repository) ItemFactory {
	return &ItemFactoryWithChildrenConcrete{
		itemRepo:    itemRepo,
		itemSetRepo: itemSetRepo,
	}
}

// CreateItem implements ItemFactory.
func (i *ItemFactoryWithChildrenConcrete) CreateItem() models.ItemInterface {
	var item models.ItemWithChildren
	return &item

}
