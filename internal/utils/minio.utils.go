package utils

import (
	"strings"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
)

func ExtractUrl(url string) (string, string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", "", exceptions.ErrInvalidS3Url
	}

	return parts[0], parts[1], nil
}
