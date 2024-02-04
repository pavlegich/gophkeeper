package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	errs "github.com/pavlegich/gophkeeper/internal/client/errors"
)

// IsValidDataType checks whether the data type is correct.
func IsValidDataType(t string) bool {
	if t == "credentials" || t == "text" || t == "binary" || t == "card" {
		return true
	}
	return false
}

// IsValidCardNumber checks whether the bank card number
// is valid using Luhn algorithm.
func IsValidCardNumber(number int) bool {
	if len(strconv.Itoa(number)) != 16 {
		return false
	}
	return (number%10+checksum(number/10))%10 == 0
}

// checkSum checks one part of bank card number
// for validity using Luhn algorithm.
func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}

func SaveFromMultipartToFile(ctx context.Context, r *http.Response, path string) error {
	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return fmt.Errorf("SaveFromMultipartToFile: couldn't get media type %w", err)
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		return fmt.Errorf("SaveFromMultipartToFile: %w", errs.ErrNotExist)
	}

	multipartReader := multipart.NewReader(r.Body, params["boundary"])
	defer r.Body.Close()

	field, err := multipartReader.NextPart()
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("SaveFromMultipartToFile: get next multi part failed %w", err)
	}
	defer field.Close()

	if field.FormName() != "file" {
		return fmt.Errorf("SaveFromMultipartToFile: no field with name file")
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("SaveFromMultipartToFile: create file failed %w", err)
	}
	defer file.Close()
	bytes, err := io.ReadAll(field)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("SaveFromMultipartToFile: couldn't read multiform field %w", err)
	}
	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("SaveFromMultipartToFile: write to file failed %w", err)
	}

	return nil
}
