package strategy

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
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

func NewExcelExportStrategy(requestRepo RequestRepo.Repository, itemRepo ItemRequestRepo.Repository, minio Minio.Repository) ExportStrategy {
	return &ExcelExportStrategy{
		requestRepo: requestRepo,
		itemRepo:    itemRepo,
		minio:       minio,
	}
}

func (e *ExcelExportStrategy) Export(ctx context.Context, requests []models.Request) ([]byte, error) {
	var items []models.ItemRequested
	for _, request := range requests {
		if request.RequestType != enums.RequestTypeRequest {
			return nil, exceptions.ErrRequestNotSupportedExportType
		}
		itemRequested, err := e.itemRepo.FindByID(ctx, request.ItemID)
		if err != nil {
			log.Error().Err(err).Msg("failed to find item requested by request ID")
			return nil, err
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
		return nil, err
	}
	f.SetActiveSheet(index)
	err = e.setHeaders(f, sheetName)
	if err != nil {
		log.Error().Err(err).Msg("failed to set headers in excel")
		return nil, err
	}

	for i, item := range items {
		request := requests[i]
		if err := e.createRow(ctx, f, sheetName, i+2, &item, &request); err != nil {
			log.Error().Err(err).Msg("failed to create row in excel")
			return nil, err
		}
	}

	// Write to buffer instead of file
	buffer, err := f.WriteToBuffer()
	if err != nil {
		log.Error().Err(err).Msg("failed to write excel file to buffer")
		return nil, err
	}

	return buffer.Bytes(), nil
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
	err = f.SetCellValue(sheetName, "E1", "Item Quantity")
	if err != nil {
		log.Error().Err(err).Msg("failed to set cell value for header")
		return err
	}
	return nil
}

func (e *ExcelExportStrategy) createRow(ctx context.Context, f *excelize.File, sheetName string, row int, item *models.ItemRequested, request *models.Request) error {
	// Try to add image, but don't fail if it doesn't work
	if request.RequestImageURL != nil && *request.RequestImageURL != "" {
		bucket, obj, err := utils.ExtractUrl(*request.RequestImageURL)
		if err != nil {
			log.Warn().Err(err).Msg("failed to extract URL, skipping image")
		} else {
			imageData, contentType, err := e.minio.GetImage(ctx, bucket, obj)
			if err != nil {
				log.Warn().Err(err).Msg("failed to get image from minio, skipping image")
			} else {
				// Validate the image data can be decoded
				_, format, err := image.DecodeConfig(bytes.NewReader(imageData))
				if err != nil {
					log.Warn().Err(err).Str("contentType", contentType).Int("dataSize", len(imageData)).Msg("failed to decode image config, image data may be corrupted")
				} else {
					log.Info().Str("detectedFormat", format).Str("contentType", contentType).Msg("Image format detected")

					// Map format to extension
					var ext string
					switch format {
					case "jpeg", "jpg":
						ext = ".jpg"
					case "png":
						ext = ".png"
					case "gif":
						ext = ".gif"
					case "bmp":
						ext = ".bmp"
					case "tiff":
						ext = ".tiff"
					default:
						log.Warn().Str("format", format).Msg("Unsupported image format")
						ext = ""
					}

					if ext != "" {
						log.Info().Str("extension", ext).Int("imageSize", len(imageData)).Msg("Adding image to excel")

						err = f.AddPictureFromBytes(sheetName, fmt.Sprintf("A%d", row), &excelize.Picture{
							Extension: ext,
							File:      imageData,
							Format:    &excelize.GraphicOptions{ScaleX: 0.3, ScaleY: 0.3},
						})
						if err != nil {
							log.Warn().Err(err).Str("extension", ext).Str("format", format).Msg("failed to add picture to excel, skipping image")
						} else {
							log.Info().Int("row", row).Str("format", format).Msg("Successfully added image to excel")
						}
					}
				}
			}
		}
	}

	err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), item.Name)
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
