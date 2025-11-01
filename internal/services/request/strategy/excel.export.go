package strategy

import (
	"context"
	"fmt"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	ItemRequestRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_requested"
	Minio "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/minio"
	RequestRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/request"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

type ExcelExportStrategy struct {
	requestRepo RequestRepo.Repository
	itemRepo    ItemRequestRepo.Repository
	minio       Minio.Repository
}

func NewExcelExportStrategy(requestRepo RequestRepo.Repository) ExportStrategy {
	return &ExcelExportStrategy{
		requestRepo: requestRepo,
	}
}

func (e *ExcelExportStrategy) Export(ctx context.Context, requests []models.Request) error {
	var items []models.ItemRequested
	for _, request := range requests {
		itemRequested, err := e.itemRepo.FindByID(ctx, request.RequestID)
		if err != nil {
			log.Error().Err(err).Msg("failed to find item requested by request ID")
			return err
		}
		if itemRequested == nil {
			log.Error().Msg("item requested not found for request ID: " + request.RequestID.String())
			continue
		}
		items = append(items, *itemRequested)
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close excel file")
		}
	}()
	sheetName := "Requests"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		log.Error().Err(err).Msg("failed to create new sheet")
		return err
	}
	f.SetActiveSheet(index)
	err = e.setHeaders(f, sheetName)
	if err != nil {
		log.Error().Err(err).Msg("failed to set headers in excel")
		return err
	}

	for i, item := range items {
		request := requests[i]
		if err := e.createRow(ctx, f, sheetName, i+2, &item, &request); err != nil {
			log.Error().Err(err).Msg("failed to create row in excel")
			return err
		}
	}

	if err := utils.SaveAtDownload(f, "requests.xlsx"); err != nil {
		log.Error().Err(err).Msg("failed to save excel file")
		return err
	}

	return nil

}

func (e *ExcelExportStrategy) setHeaders(f *excelize.File, sheetName string) error {
	err := f.SetCellValue(sheetName, "A1", "Item Image")
	if err != nil {
		log.Error().Err(err).Msg("failed to set cell value for header")
		return err
	}
	err = f.SetCellValue(sheetName, "B1", "Item Name")
	if err != nil {
		log.Error().Err(err).Msg("failed to set cell value for header")
		return err
	}
	err = f.SetCellValue(sheetName, "C1", "Item Description")
	if err != nil {
		log.Error().Err(err).Msg("failed to set cell value for header")
		return err
	}
	err = f.SetCellValue(sheetName, "D1", "Item Type")
	if err != nil {
		log.Error().Err(err).Msg("failed to set cell value for header")
		return err
	}
	f.SetCellValue(sheetName, "E1", "Item Quantity")
	if err != nil {
		log.Error().Err(err).Msg("failed to set cell value for header")
		return err
	}
	return nil
}

func (e *ExcelExportStrategy) createRow(ctx context.Context, f *excelize.File, sheetName string, row int, item *models.ItemRequested, request *models.Request) error {
	bucket, obj, err := utils.ExtractUrl(*request.RequestImageURL)
	if err != nil {
		log.Error().Err(err).Msg("failed to extract URL")
		return err
	}
	image, objType, err := e.minio.GetImage(ctx, bucket, obj)
	if err != nil {
		log.Error().Err(err).Msg("failed to get image from minio")
		return err
	}

	err = f.AddPictureFromBytes(sheetName, fmt.Sprintf("A%d", row), &excelize.Picture{
		Extension: objType,
		File:      image,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to add picture to excel")
		return err
	}

	err = f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), item.Name)
	if err != nil {
		log.Error().Err(err).Msg("failed to set cell value for item name")
		return err
	}
	err = f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), item.Description)
	if err != nil {
		log.Error().Err(err).Msg("failed to set cell value for item description")
		return err
	}
	err = f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), item.Type)
	if err != nil {
		log.Error().Err(err).Msg("failed to set cell value for item type")
		return err
	}
	err = f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), item.Quantity)
	if err != nil {
		log.Error().Err(err).Msg("failed to set cell value for item quantity")
		return err
	}

	return nil
}
