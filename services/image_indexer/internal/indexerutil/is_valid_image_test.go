package indexerutil

import (
	"bytes"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createTestImage(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, img)
	if err != nil {
		return nil
	}

	return buffer.Bytes()
}

func TestIsValidImage(t *testing.T) {
	tests := []struct {
		name        string
		imageBytes  []byte
		statusCode  int
		minWidth    int
		minHeight   int
		expectValid bool
	}{
		{
			name:        "valid image meets size",
			imageBytes:  createTestImage(200, 500),
			statusCode:  http.StatusOK,
			minWidth:    ImgMinWidth,
			minHeight:   ImgMinHeight,
			expectValid: true,
		},
		{
			name:        "image too small",
			imageBytes:  createTestImage(50, 50),
			statusCode:  http.StatusOK,
			minWidth:    ImgMinWidth,
			minHeight:   ImgMinHeight,
			expectValid: false,
		},
		{
			name:        "non-200 response",
			imageBytes:  createTestImage(200, 200),
			statusCode:  http.StatusNotFound,
			minWidth:    ImgMinWidth,
			minHeight:   ImgMinHeight,
			expectValid: false,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				if tc.statusCode == http.StatusOK {
					w.Write(tc.imageBytes)
				}
			}))
			defer testServer.Close()

			actual := IsValidImage(testServer.URL, tc.minWidth, tc.minHeight)
			if actual != tc.expectValid {
				t.Errorf("Test %d - '%s' FAIL: expected %v, got %v", i, tc.name, tc.expectValid, actual)
			}
		})
	}
}
