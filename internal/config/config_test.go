// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package config

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/alpstable/gidari"
	"golang.org/x/time/rate"
)

func TestReadFile(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name string
		path string
		want *gidari.Config
	}{
		{
			name: "valid",
			path: "testdata/valid.yaml",
			want: &gidari.Config{
				RawURL: "https://chroniclingamerica.loc.gov",
				RateLimitConfig: &gidari.RateLimitConfig{
					Burst:  intPtr(1),
					Period: timeDurPtr(time.Second),
				},
				Requests: []*gidari.Request{
					{
						Endpoint: "/search/titles/results",
						Query: map[string]string{
							"terms":  "oakland",
							"format": "json",
							"page":   "5",
						},
					},
				},
				StorageOptions: []gidari.StorageOptions{
					{
						ConnectionString: strPtr("mongodb://mongo1:27017"),
						Database:         strPtr("histam"),
					},
				},
			},
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			got, err := readFile(tcase.path)
			if err != nil {
				t.Fatalf("error reading file: %v", err)
			}

			if got.RawURL != tcase.want.RawURL {
				t.Fatalf("got RawURL %s, want %s", got.RawURL, tcase.want.RawURL)
			}

			if got.RateLimitConfig == nil {
				t.Fatal("got nil RateLimitConfig")
			}

			if got.RateLimitConfig.Burst == nil {
				t.Fatal("got nil RateLimitConfig.Burst")
			}

			if got.RateLimitConfig.Period == nil {
				t.Fatal("got nil RateLimitConfig.Period")
			}

			if *got.RateLimitConfig.Burst != *tcase.want.RateLimitConfig.Burst {
				t.Fatalf("got RateLimitConfig.Burst %d, want %d", *got.RateLimitConfig.Burst,
					*tcase.want.RateLimitConfig.Burst)
			}

			if *got.RateLimitConfig.Period != *tcase.want.RateLimitConfig.Period {
				t.Fatalf("got RateLimitConfig.Period %s, want %s", *got.RateLimitConfig.Period,
					*tcase.want.RateLimitConfig.Period)
			}

			if len(got.Requests) != len(tcase.want.Requests) {
				t.Fatalf("got %d requests, want %d", len(got.Requests), len(tcase.want.Requests))
			}

			for idx, req := range got.Requests {
				if req.Endpoint != tcase.want.Requests[idx].Endpoint {
					t.Fatalf("got request %d endpoint %s, want %s", idx, req.Endpoint,
						tcase.want.Requests[idx].Endpoint)
				}

				if len(req.Query) != len(tcase.want.Requests[idx].Query) {
					t.Fatalf("got %d query params, want %d", len(req.Query),
						len(tcase.want.Requests[idx].Query))
				}

				for k, v := range req.Query {
					if v != tcase.want.Requests[idx].Query[k] {
						t.Fatalf("got query param %s=%s, want %s=%s", k, v, k,
							tcase.want.Requests[idx].Query[k])
					}
				}
			}

			if len(got.StorageOptions) != len(tcase.want.StorageOptions) {
				t.Fatalf("got %d storage options, want %d", len(got.StorageOptions),
					len(tcase.want.StorageOptions))
			}

			for _, sto := range got.StorageOptions {
				if sto.ConnectionString == nil {
					t.Fatal("got nil StorageOptions.ConnectionString")
				}

				if sto.Database == nil {
					t.Fatal("got nil StorageOptions.Database")
				}
			}
		})
	}
}

func TestAllStorage(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name   string
		config *gidari.Config
	}{
		{
			name: "valid",
			config: &gidari.Config{
				StorageOptions: []gidari.StorageOptions{
					{
						ConnectionString: strPtr("mongodb://mongo1:27017"),
						Database:         strPtr("test"),
					},
				},
			},
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			var err error

			tcase.config.StorageOptions, err = addAllStorage(context.Background(),
				tcase.config.StorageOptions)

			if err != nil {
				t.Fatalf("error getting storage: %v", err)
			}

			for _, sto := range tcase.config.StorageOptions {
				stg := sto.Storage
				if stg == nil {
					t.Fatal("got nil storage.Storage")
				}

				if err := stg.Ping(); err != nil {
					t.Fatalf("error pinging storage: %v", err)
				}

				stg.Close()
			}
		})
	}
}

func TestAddRequestData(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name   string
		config *gidari.Config
		want   []*gidari.Request
	}{
		{
			name: "valid",
			config: &gidari.Config{
				RateLimitConfig: &gidari.RateLimitConfig{
					Burst:  intPtr(10),
					Period: timeDurPtr(1 * time.Second),
				},
				Requests: []*gidari.Request{
					{
						Endpoint: "/repos/golang/go/issues",
						Query: map[string]string{
							"format": "json",
							"page":   "5",
						},
					},
				},
			},
			want: []*gidari.Request{
				{
					Endpoint: "/repos/golang/go/issues",
					Query: map[string]string{
						"format": "json",
						"page":   "5",
					},
					RateLimiter: rate.NewLimiter(rate.Every(1*time.Second), 10),
					Method:      http.MethodGet,
					Table:       "issues",
				},
			},
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			rlc := tcase.config.RateLimitConfig
			addRequestData(context.Background(), rlc, tcase.config.Requests)

			for idx, req := range tcase.config.Requests {
				if req.RateLimiter == nil {
					t.Fatal("got nil RateLimiter")
				}

				if req.RateLimiter.Limit() != tcase.want[idx].RateLimiter.Limit() {
					t.Fatalf("got rate limit %f, want %f", req.RateLimiter.Limit(),
						tcase.want[idx].RateLimiter.Limit())
				}

				if req.RateLimiter.Burst() != tcase.want[idx].RateLimiter.Burst() {
					t.Fatalf("got rate burst %d, want %d", req.RateLimiter.Burst(),
						tcase.want[idx].RateLimiter.Burst())
				}

				if req.Method == "" {
					t.Fatal("got empty Method")
				}

				if req.Method != tcase.want[idx].Method {
					t.Fatalf("got method %s, want %s", req.Method, tcase.want[idx].Method)
				}

				if req.Table == "" {
					t.Fatal("got empty Table")
				}

				if req.Table != tcase.want[idx].Table {
					t.Fatalf("got table %s, want %s", req.Table, tcase.want[idx].Table)
				}

				if req.Endpoint == "" {
					t.Fatal("got empty Endpoint")
				}

				if req.Endpoint != tcase.want[idx].Endpoint {
					t.Fatalf("got endpoint %s, want %s", req.Endpoint, tcase.want[idx].Endpoint)
				}
			}
		})
	}
}

// intPtr will convert an int to a pointer to an int.
func intPtr(i int) *int {
	return &i
}

// strPtr will conver a string to a pointer to a string.
func strPtr(str string) *string {
	return &str
}

func timeDurPtr(d time.Duration) *time.Duration {
	return &d
}
