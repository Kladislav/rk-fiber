// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package rkfiberauth

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	rkfiberinter "github.com/rookie-ninja/rk-fiber/interceptor"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInterceptor_WithIgnoringPath(t *testing.T) {
	defer assertNotPanic(t)

	app := fiber.New()

	handler := Interceptor(
		WithEntryNameAndType("ut-entry", "ut-type"),
		WithBasicAuth("ut-realm", "user:pass"),
		WithApiKeyAuth("ut-api-key"),
		WithIgnorePrefix("/ut-ignore-path"))

	app.Use(handler)
	app.Get("/ut-ignore-path", func(ctx *fiber.Ctx) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/ut-ignore-path", nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestInterceptor_WithBasicAuth_Invalid(t *testing.T) {
	defer assertNotPanic(t)

	app := fiber.New()

	handler := Interceptor(
		WithEntryNameAndType("ut-entry", "ut-type"),
		WithBasicAuth("ut-realm", "user:pass"))

	app.Use(handler)
	app.Get("/ut-path", func(ctx *fiber.Ctx) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/ut-path", nil)
	req.Header.Set(rkfiberinter.RpcAuthorizationHeaderKey, "invalid")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestInterceptor_WithBasicAuth_InvalidBasicAuth(t *testing.T) {
	defer assertNotPanic(t)

	app := fiber.New()

	handler := Interceptor(
		WithEntryNameAndType("ut-entry", "ut-type"),
		WithBasicAuth("ut-realm", "user:pass"))

	app.Use(handler)
	app.Get("/ut-path", func(ctx *fiber.Ctx) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/ut-path", nil)
	req.Header.Set(rkfiberinter.RpcAuthorizationHeaderKey, fmt.Sprintf("%s invalid", typeBasic))
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestInterceptor_WithApiKey_Invalid(t *testing.T) {
	defer assertNotPanic(t)

	app := fiber.New()

	handler := Interceptor(
		WithEntryNameAndType("ut-entry", "ut-type"),
		WithApiKeyAuth("ut-api-key"))

	app.Use(handler)
	app.Get("/ut-path", func(ctx *fiber.Ctx) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/ut-path", nil)
	req.Header.Set(rkfiberinter.RpcApiKeyHeaderKey, "invalid")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestInterceptor_MissingAuth(t *testing.T) {
	defer assertNotPanic(t)

	app := fiber.New()

	handler := Interceptor(
		WithEntryNameAndType("ut-entry", "ut-type"),
		WithApiKeyAuth("ut-api-key"))

	app.Use(handler)
	app.Get("/ut-path", func(ctx *fiber.Ctx) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/ut-path", nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestInterceptor_HappyCase(t *testing.T) {
	defer assertNotPanic(t)

	app := fiber.New()

	handler := Interceptor(
		WithEntryNameAndType("ut-entry", "ut-type"),
		//WithBasicAuth("ut-realm", "user:pass"),
		WithApiKeyAuth("ut-api-key"))

	app.Use(handler)
	app.Get("/ut-path", func(ctx *fiber.Ctx) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/ut-path", nil)
	req.Header.Set(rkfiberinter.RpcApiKeyHeaderKey, "ut-api-key")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func assertNotPanic(t *testing.T) {
	if r := recover(); r != nil {
		// Expect panic to be called with non nil error
		assert.True(t, false)
	} else {
		// This should never be called in case of a bug
		assert.True(t, true)
	}
}