# Changelog

## v3.0.0 (2019-04-15)

### Features

* Add options for : `DisableInternalErrorMiddleware`, `DisableRecoveryMiddleware`, `DisablePingRouter` [PR](https://github.com/ifreddyrondon/bastion/pull/13)
* Expose logger request as a middleware with their config functions. [PR](https://github.com/ifreddyrondon/bastion/pull/14)
    * `AttachLogger(log zerolog.Logger)` chain the logger with the middleware.
    * `EnableLogReqIP()` show the request ip.
    * `EnableLogUserAgent()` show the user agent of the request.
    * `EnableLogReferer()` show referer of the request.
    * `DisableLogMethod()` hide the request method.
    * `DisableLogURL()` hide the request url.
    * `DisableLogStatus()` hide the request status.
    * `DisableLogSize()` hide the request size.
    * `DisableLogDuration()` hide the request duration.
    * `DisableLogRequestID()` hide the request id.
* Add option to disable logger middleware `DisableLoggerMiddleware`. [PR](https://github.com/ifreddyrondon/bastion/pull/14)
* Add debug mode. [PR](https://github.com/ifreddyrondon/bastion/pull/15)
* Add default JSON renderer `render.JSON` [PR](https://github.com/ifreddyrondon/bastion/pull/16)
* Add handler for not found resource. [PR](https://github.com/ifreddyrondon/bastion/pull/17)
* Add handler for not allowed method in resource. [PR](https://github.com/ifreddyrondon/bastion/pull/17)
* Add /debug router used for mounting `net/http/pprof`. [PR](https://github.com/ifreddyrondon/bastion/pull/18)

### Breaking changes
* Remove internal router and APIRouter attribute, now the Bastion instance is the mux. _breaking changes_ [PR](https://github.com/ifreddyrondon/bastion/pull/8)
* Remove support for load config from file. _breaking changes_ [PR](https://github.com/ifreddyrondon/bastion/pull/9)
* Remove address attribute from Config and pass it when `Serve` or use the default. It also can be set with ADDR env. _breaking changes_. [PR](https://github.com/ifreddyrondon/bastion/pull/10)
* Bastion logger is now internal. It can be used from context. _breaking changes_. [PR](https://github.com/ifreddyrondon/bastion/pull/11) 
* Rename `APIError` middleware including bastion option `API500ErrMessage ` to `InternalError` and `InternalErrMsg`. _breaking changes_. [PR](https://github.com/ifreddyrondon/bastion/pull/12).
* Rename `NoPrettyLogging` option to `DisablePrettyLogging`. [PR](https://github.com/ifreddyrondon/bastion/pull/13)

- - - 
* History of changes: see https://github.com/ifreddyrondon/bastion/compare/v2.3.0...v3.0.0

## v2.3.0 (2019-04-05)

- Minor release
- Update github.com/go-chi/chi to v4.0.2
- Update github.com/pkg/errors to 0.8.1
- Update github.com/rs/zerolog to 1.13.0
- Remove github.com/gobuffalo/envy dependency
- Remove github.com/markbates/going dependency
- History of changes: see https://github.com/ifreddyrondon/bastion/compare/v2.2.0...v2.3.0