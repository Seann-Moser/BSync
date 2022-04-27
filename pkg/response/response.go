package response

import (
	"encoding/json"
	"github.com/Seann-Moser/BaseGoAPI/pkg/pagination"
	"go.uber.org/zap"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Response struct {
	logger *zap.Logger
}
type BaseResponse struct {
	Message string                 `json:"message"`
	Data    interface{}            `json:"data,omitempty"`
	Page    *pagination.Pagination `json:"page,omitempty"`
}

func NewResponse(logger *zap.Logger) *Response {
	return &Response{logger: logger}
}

func (resp *Response) Error(w http.ResponseWriter, err error, code int, message string) {
	w.WriteHeader(code)
	EncodeErr := json.NewEncoder(w).Encode(BaseResponse{
		Message: message,
	})
	if EncodeErr != nil {
		resp.logger.Error("failed encoding response", zap.Error(EncodeErr))
	}
}

func (resp *Response) PaginationResponse(w http.ResponseWriter, data []interface{}, page *pagination.Pagination) {
	w.WriteHeader(http.StatusOK)
	bytes, err := json.MarshalIndent(BaseResponse{
		Data: getRange(data, page),
		Page: page,
	}, "", "    ")
	if err != nil {
		resp.logger.Error("failed to encode response")
	}
	_, EncodeErr := w.Write(bytes)
	if EncodeErr != nil {
		resp.logger.Error("failed encoding response", zap.Error(EncodeErr))
	}
}
func getRange(data []interface{}, page *pagination.Pagination) []interface{} {
	page.TotalItems = uint(len(data))
	page.TotalPages = uint(math.Ceil(float64(page.TotalItems) / float64(page.ItemsPerPage)))
	if len(data) < int(page.ItemsPerPage) {
		return data
	}

	min := int(page.CurrentPage * page.ItemsPerPage)
	max := min + int(page.ItemsPerPage)

	if len(data) > min {
		return nil
	}
	if len(data) < max {
		return data[min:]
	}
	return data[min:max]
}

func (resp *Response) Message(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusOK)
	bytes, err := json.MarshalIndent(BaseResponse{
		Message: msg,
	}, "", "    ")
	if err != nil {
		resp.logger.Error("failed to encode response")
	}
	_, EncodeErr := w.Write(bytes)
	if EncodeErr != nil {
		resp.logger.Error("failed encoding response", zap.Error(EncodeErr))
	}
}

func (resp *Response) DataResponse(w http.ResponseWriter, data interface{}, code int) {
	w.WriteHeader(code)
	bytes, err := json.MarshalIndent(BaseResponse{
		Data: data,
	}, "", "    ")
	if err != nil {
		resp.logger.Error("failed to encode response")
	}
	_, EncodeErr := w.Write(bytes)
	if EncodeErr != nil {
		resp.logger.Error("failed encoding response", zap.Error(EncodeErr))
	}
}

func (resp *Response) File(w http.ResponseWriter, file string, download bool) (int64, error) {
	filename := strings.Split(file, "/")
	w.Header().Set("filename", filename[len(filename)-1])
	if download {
		w.Header().Set("Content-Description", "File Transfer")
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filename[len(filename)-1]))
		//w.Header().Set("Content-Type", "application/octet-stream")
	}
	f, _ := os.Open(file)
	defer func() {
		_ = f.Close()
	}()

	fileHeader := make([]byte, 512)
	_, err := f.Read(fileHeader)
	if err != nil {
		return 0, err
	}
	fileStat, _ := f.Stat()
	w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	_, err = f.Seek(0, 0)
	if err != nil {
		return 0, err
	}
	return io.Copy(w, f)
}
