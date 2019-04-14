package middleware

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"
	"github.com/ifreddyrondon/bastion/render"
)

var (
	// ListingCtxKey is the context.Context key to store the Listing for a request.
	ListingCtxKey = &contextKey{"Listing"}
)

var (
	errMissingListing    = errors.New("listing not found in context")
	errWrongListingValue = errors.New("listing value set incorrectly in context")
)

func withParams(ctx context.Context, l *listing.Listing) context.Context {
	return context.WithValue(ctx, ListingCtxKey, l)
}

// GetListing will return the listing reference assigned to the context, or nil if there
// is any error or there isn't a Listing instance.
func GetListing(ctx context.Context) (*listing.Listing, error) {
	tmp := ctx.Value(ListingCtxKey)
	if tmp == nil {
		return nil, errMissingListing
	}
	l, ok := tmp.(*listing.Listing)
	if !ok {
		return nil, errWrongListingValue
	}
	return l, nil
}

type listingConfig struct {
	render         render.ClientErrRenderer
	optionsDecoder []func(*listing.Decoder)
}

func getListingCfg(opts ...func(*listingConfig)) *listingConfig {
	r := &listingConfig{
		render: render.JSON,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// Limit set the paging limit default.
func Limit(limit int) func(*listingConfig) {
	return func(l *listingConfig) {
		o := listing.DecodeLimit(limit)
		l.optionsDecoder = append(l.optionsDecoder, o)
	}
}

// MaxAllowedLimit set the max allowed limit default.
func MaxAllowedLimit(maxAllowed int) func(*listingConfig) {
	return func(l *listingConfig) {
		o := listing.DecodeMaxAllowedLimit(maxAllowed)
		l.optionsDecoder = append(l.optionsDecoder, o)
	}
}

// Sort set criteria to sort
func Sort(criteria ...sorting.Sort) func(*listingConfig) {
	return func(l *listingConfig) {
		o := listing.DecodeSort(criteria...)
		l.optionsDecoder = append(l.optionsDecoder, o)
	}
}

// Filter set criteria to filter
func Filter(criteria ...filtering.FilterDecoder) func(*listingConfig) {
	return func(l *listingConfig) {
		o := listing.DecodeFilter(criteria...)
		l.optionsDecoder = append(l.optionsDecoder, o)
	}
}

// Listing is a middleware that parses the url from a request and stores a
// listing.Listing on the context, it can be accessed through middleware.GetListing.
//
// Sample usage.. for the url: `/repositories/1?limit=10&offset=25`
//
//  func routes() http.Handler {
//    r := chi.NewRouter()
//    r.Use(middleware.Listing())
//    r.Get("/repositories/{id}", ListRepositories)
//    return r
//  }
//
//  func ListRepositories(w http.ResponseWriter, r *http.Request) {
// 	  list, _ := middleware.GetListing(r.Context())
//
// 	  // do something with listing
// }
func Listing(opts ...func(*listingConfig)) func(http.Handler) http.Handler {
	cfg := getListingCfg(opts...)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var l listing.Listing

			if err := listing.NewDecoder(r.URL.Query(), cfg.optionsDecoder...).Decode(&l); err != nil {
				cfg.render.BadRequest(w, err)
				return
			}

			ctx := withParams(r.Context(), &l)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
