package helper

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGetImageFileFromUrl(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sample image content"))
	}))
	defer server.Close()

	imageReader := GetImageFileFromUrl(server.URL)

	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, imageReader)
	if err != nil {
		t.Errorf("Error reading image content: %v", err)
	}

	expectedContent := "sample image content"
	if buffer.String() != expectedContent {
		t.Errorf("Expected content: %s, got: %s", expectedContent, buffer.String())
	}
}

func TestGetImageFileFromPath(t *testing.T) {
	fnOSOpen = func(name string) (*os.File, error) {
		// seed mockFile with mock content
		mockFile := &os.File{}

		return mockFile, nil
	}

	fnIOCopy = func(dst io.Writer, src io.Reader) (written int64, err error) {
		mockContent := "sample content"
		io.Copy(dst, strings.NewReader(mockContent))
		return 1, nil
	}

	existingImagePath := "/path/to/your/existing/image.jpg"

	imageReader := GetImageFileFromPath(existingImagePath)

	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, imageReader)
	if err != nil {
		t.Errorf("Error reading image content: %v", err)
	}

	if buffer.Len() == 0 {
		t.Error("Empty image content")
	}
}
