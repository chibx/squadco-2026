package routes

import (
	"mime/multipart"

	server "github.com/chibx/vendor-pulse/internal/server_errors"
)

type myFile struct {
	Size int64
	F    multipart.File
}

var businessDocs = []string{"cac_document", "nin_document", "logo", "face"}

func allFilesExist(form *multipart.Form, names []string) (map[string]myFile, []*server.ErrorDetail, error) {
	fileMap := make(map[string]myFile)
	errorBag := make([]*server.ErrorDetail, 0, len(names))

	for _, v := range names {
		files, ok := form.File["cac_document"]
		if !ok || len(files) == 0 {
			errorBag = append(errorBag, businesDocsErrorDetail(v))
			continue
		}

		file, err := files[0].Open()
		if err != nil {
			return nil, nil, err
		}
		fileMap[v] = myFile{
			Size: files[0].Size,
			F:    file,
		}
	}

	return fileMap, errorBag, nil
}

func businesDocsErrorDetail(name string) *server.ErrorDetail {
	switch name {
	case "cac_document":
		return &server.ErrorDetail{
			Field:   "CAC File",
			Message: "CAC document is missing",
		}
	case "nin_document":
		return &server.ErrorDetail{
			Field:   "NIN File",
			Message: "NIN document is missing",
		}
	case "logo":
		return &server.ErrorDetail{
			Field:   "Business Logo",
			Message: "Brand Logo is missing",
		}
	default:
		return &server.ErrorDetail{
			Field:   "Face Photo",
			Message: "Portrait is missing",
		}
	}
}
