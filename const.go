// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

// HTTPMethod is a set of request methods to indicate the purpose of the request
// and what is expected if the request is successful.
//
// Each request method has its own semantics, but some characteristics
// are shared across multiple methods, specifically request methods
// can be [safe], [idempotent], or [cacheable].
//
// [safe]: https://developer.mozilla.org/en-US/docs/Glossary/Safe/HTTP
// [idempotent]: https://developer.mozilla.org/en-US/docs/Glossary/Idempotent
// [cacheable]: https://developer.mozilla.org/en-US/docs/Glossary/Cacheable
type HTTPMethod string

const (
	// The GET method requests a representation of the specified resource.
	// Requests using GET should only retrieve data and should not contain
	// a request content.
	//
	// Semantics:
	//  - Safe ✅
	//  - Idempotent ✅
	//  - Cacheable ✅
	GET HTTPMethod = "GET"

	// The POST method submits an entity to the specified resource,
	// often causing a change in state or side effects on the server.
	//
	// Semantics:
	//  - Safe ❌
	//  - Idempotent ❌
	//  - Cacheable when responses explicitly include freshness information
	//    and a matching Content-Location header.
	POST HTTPMethod = "POST"

	// The PUT method replaces all current representations of the target
	// resource with the request content.
	//
	// Semantics:
	//  - Safe ❌
	//  - Idempotent ✅
	//  - Cacheable ❌
	PUT HTTPMethod = "PUT"

	// The DELETE method deletes the specified resource.
	//
	// Semantics:
	//  - Safe ❌
	//  - Idempotent ✅
	//  - Cacheable ❌
	DELETE HTTPMethod = "DELETE"

	// The OPTIONS method describes the communication options for the target
	// resource.
	//
	// Semantics:
	//  - Safe ✅
	//  - Idempotent ✅
	//  - Cacheable ❌
	OPTIONS HTTPMethod = "OPTIONS"

	// The PATCH method applies partial modifications to a resource.
	//
	// Semantics:
	//  - Safe ❌
	//  - Idempotent ❌
	//  - Cacheable when responses explicitly include freshness information
	//    and a matching Content-Location header.
	PATCH HTTPMethod = "PATCH"
)
