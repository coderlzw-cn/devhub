package response

import (
	"errors"
	"net/http"
	"testing"
)

func TestServerErrorMessage(t *testing.T) {
	err := NewError(nil, http.StatusInternalServerError, "Internal server error", errors.New("database details"))
	if err.Error() != "服务暂时不可用，请稍后重试" {
		t.Fatalf("unexpected server error message: %s", err.Error())
	}
	clientErr := NewError(nil, http.StatusBadRequest, "参数错误", errors.New("分类 ID 格式不正确"))
	if clientErr.Error() != "分类 ID 格式不正确" {
		t.Fatalf("lost client error details: %s", clientErr.Error())
	}
}
