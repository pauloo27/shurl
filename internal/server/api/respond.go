package api

import (
	"net/http"
)

func Ok[T any](detail T) Response {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	return Response{
		StatusCode: http.StatusOK,
		Body:       detail,
		Header:     header,
	}
}

func Created[T any](detail T) Response {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	return Response{
		StatusCode: http.StatusCreated,
		Body:       detail,
		Header:     header,
	}
}

func Redirect(url string) Response {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("Location", url)

	return Response{
		StatusCode: http.StatusTemporaryRedirect,
		Header:     header,
	}
}

func Err(
	error ErrorType,
	message string,
) Response {
	return DetailedError(error, map[string]string{
		"message": message,
	})
}

func DetailedError[T any](
	error ErrorType,
	detail T,
) Response {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	return Response{
		StatusCode: error.StatusCode,
		Body: Error[T]{
			Error:  error.Name,
			Detail: detail,
		},
		Header: header,
	}
}
