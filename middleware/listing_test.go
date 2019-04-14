package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gavv/httpexpect.v1"

	"github.com/ifreddyrondon/bastion/middleware"
	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"
	"github.com/ifreddyrondon/bastion/middleware/listing/paging"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"
	"github.com/ifreddyrondon/bastion/render"
)

func TestGetListingMissingInstance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	_, err := middleware.GetListing(ctx)
	assert.EqualError(t, err, "listing not found in context")
}

func TestGetListingInvalidReference(t *testing.T) {
	t.Parallel()
	ctx := context.WithValue(context.Background(), middleware.ListingCtxKey, 1)
	_, err := middleware.GetListing(ctx)
	assert.EqualError(t, err, "listing value set incorrectly in context")
}

func setup(m func(http.Handler) http.Handler) (*httptest.Server, *listing.Listing, func()) {
	var result listing.Listing

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l, err := middleware.GetListing(r.Context())
		if err != nil {
			render.JSON.InternalServerError(w, err)
			return
		}
		result = *l
		w.Write([]byte("hi"))
	})

	server := httptest.NewServer(m(h))
	teardown := func() {
		server.Close()
	}
	return server, &result, teardown
}

func TestListingMiddlewareFailure(t *testing.T) {
	t.Parallel()

	createdDESC := sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")

	tt := []struct {
		name      string
		urlParams string
		m         func(http.Handler) http.Handler
		response  map[string]interface{}
	}{
		{
			"given a bad offset param should return a 400",
			"offset=abc",
			middleware.Listing(middleware.Limit(50)),
			map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "invalid offset value, must be a number",
			},
		},
		{
			"given a sort query when none match sorting criteria should return a 400",
			"sort=foo_desc",
			middleware.Listing(middleware.Sort(createdDESC)),
			map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "there's no order criteria with the id foo_desc",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server, _, teardown := setup(tc.m)
			defer teardown()
			e := httpexpect.New(t, server.URL)
			e.GET("/").WithQueryString(tc.urlParams).
				Expect().
				Status(http.StatusBadRequest).
				JSON().
				Object().Equal(tc.response)
		})
	}
}

func TestListingMiddlewareOkWithOptions(t *testing.T) {
	t.Parallel()

	createdDESC := sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")
	createdASC := sorting.NewSort("created_at_asc", "created_at ASC", "Created date ascendant")
	vNew := filtering.NewValue("new", "New")
	vUsed := filtering.NewValue("used", "Used")
	text := filtering.NewText("condition", "test", vNew, vUsed)
	vTrue := filtering.NewValue("true", "shared")
	vFalse := filtering.NewValue("false", "private")
	boolean := filtering.NewBoolean("shared", "test", "shared", "private")

	tt := []struct {
		name      string
		urlParams string
		m         func(http.Handler) http.Handler
		result    listing.Listing
	}{
		{
			"given non query params and not options should get default paging",
			"",
			middleware.Listing(),
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
				}
			}(),
		},
		{
			"given offset=11 params and not options should get paging with offset=11 and defaults",
			"offset=11",
			middleware.Listing(),
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          11,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
				}
			}(),
		},
		{
			"given non query params and changing the default limit option should get default paging with limit 50",
			"",
			middleware.Listing(middleware.Limit(50)),
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           50,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
				}
			}(),
		},
		{
			"given limit=110 param and changing the default MaxAllowedLimit option to 120 should allow limit > 100 < 120",
			"limit=110",
			middleware.Listing(middleware.MaxAllowedLimit(120)),
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           110,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: 120,
					},
				}
			}(),
		},
		{
			"given non query params and one sort criteria should get sorting with default sort",
			"",
			middleware.Listing(middleware.Sort(createdDESC)),
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Sorting: &sorting.Sorting{
						Sort:      &createdDESC,
						Available: []sorting.Sort{createdDESC},
					},
				}
			}(),
		},
		{
			"given sort query params and sort criteria should get sorting with selected sort",
			"sort=created_at_desc",
			middleware.Listing(middleware.Sort(createdDESC, createdASC)),
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Sorting: &sorting.Sorting{
						Sort:      &createdDESC,
						Available: []sorting.Sort{createdDESC, createdASC},
					},
				}
			}(),
		},
		{
			"given non query params and one filter criteria should get filtering with only available",
			"",
			middleware.Listing(middleware.Filter(text)),
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Filtering: &filtering.Filtering{
						Filters: []filtering.Filter{},
						Available: []filtering.Filter{
							{
								ID:          "condition",
								Description: "test",
								Type:        "text",
								Values:      []filtering.Value{vNew, vUsed},
							},
						},
					},
				}
			}(),
		},
		{
			"given non query params and some filters criteria should get filtering with all available",
			"",
			middleware.Listing(middleware.Filter(text, boolean)),
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Filtering: &filtering.Filtering{
						Filters: []filtering.Filter{},
						Available: []filtering.Filter{
							{
								ID:          "condition",
								Description: "test",
								Type:        "text",
								Values:      []filtering.Value{vNew, vUsed},
							},
							{
								ID:          "shared",
								Description: "test",
								Type:        "boolean",
								Values:      []filtering.Value{vTrue, vFalse},
							},
						},
					},
				}
			}(),
		},
		{
			"given a filter query params and some filters criteria should get filtering with all available and filter",
			"condition=new",
			middleware.Listing(middleware.Filter(text, boolean)),
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Filtering: &filtering.Filtering{
						Filters: []filtering.Filter{
							{
								ID:          "condition",
								Description: "test",
								Type:        "text",
								Values:      []filtering.Value{vNew},
							},
						},
						Available: []filtering.Filter{
							{
								ID:          "condition",
								Description: "test",
								Type:        "text",
								Values:      []filtering.Value{vNew, vUsed},
							},
							{
								ID:          "shared",
								Description: "test",
								Type:        "boolean",
								Values:      []filtering.Value{vTrue, vFalse},
							},
						},
					},
				}
			}(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server, resultContainer, teardown := setup(tc.m)
			defer teardown()
			e := httpexpect.New(t, server.URL)
			e.GET("/").WithQueryString(tc.urlParams).
				Expect().
				Status(http.StatusOK)

			assert.Equal(t, tc.result.Paging, resultContainer.Paging)
			if resultContainer.Sorting != nil {
				assert.Equal(t, tc.result.Sorting.Sort, resultContainer.Sorting.Sort)
				assert.Equal(t, tc.result.Sorting.Available, resultContainer.Sorting.Available)
			}
			if resultContainer.Filtering != nil {
				for i, f := range resultContainer.Filtering.Filters {
					assert.Equal(t, tc.result.Filtering.Filters[i], f)
				}
				assert.Equal(t, tc.result.Filtering.Available, resultContainer.Filtering.Available)
			}
		})
	}
}

func TestListingMiddlewareMarshall(t *testing.T) {
	t.Parallel()

	updatedDESC := sorting.NewSort("updated_at_desc", "updated_at DESC", "Updated date descending")
	updatedASC := sorting.NewSort("updated_at_asc", "updated_at ASC", "Updated date ascendant")
	createdDESC := sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")
	createdASC := sorting.NewSort("created_at_asc", "created_at ASC", "Created date ascendant")

	tt := []struct {
		name      string
		urlParams string
		m         func(http.Handler) http.Handler
		expected  string
	}{
		{
			"given non query params and sort criteria should get default sort",
			"",
			middleware.Listing(
				middleware.Sort(updatedDESC, updatedASC, createdDESC, createdASC),
			),
			`"sorting":{"sort":{"id":"updated_at_desc","description":"Updated date descending"},"available":[{"id":"updated_at_desc","description":"Updated date descending"},{"id":"updated_at_asc","description":"Updated date ascendant"},{"id":"created_at_desc","description":"Created date descending"},{"id":"created_at_asc","description":"Created date ascendant"}]}`,
		},
		{
			"given created_at_asc params and sort criteria should get created_at_asc as sort and all sort criteria",
			"sort=created_at_asc",
			middleware.Listing(
				middleware.Sort(updatedDESC, updatedASC, createdDESC, createdASC),
			),
			`"sorting":{"sort":{"id":"created_at_asc","description":"Created date ascendant"},"available":[{"id":"updated_at_desc","description":"Updated date descending"},{"id":"updated_at_asc","description":"Updated date ascendant"},{"id":"created_at_desc","description":"Created date descending"},{"id":"created_at_asc","description":"Created date ascendant"}]}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server, resultContainer, teardown := setup(tc.m)
			defer teardown()
			e := httpexpect.New(t, server.URL)
			e.GET("/").WithQueryString(tc.urlParams).
				Expect().
				Status(http.StatusOK)

			result, err := resultContainer.MarshalJSON()
			assert.Nil(t, err)
			assert.Contains(t, string(result), tc.expected)
		})
	}
}
