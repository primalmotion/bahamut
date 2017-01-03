// Author: Antoine Mercadal
// See LICENSE file for full LICENSE
// Copyright 2016 Aporeto.

package bahamut

import (
	"net/http"

	"github.com/aporeto-inc/elemental"
)

// CheckAuthentication checks if the current context has been authenticated if there is any authenticator registered.
//
// If it is not authenticated it stops the normal processing execution flow, and will write the Unauthorized response to the given writer.
// If not Authenticator is set, then it will always return true.
//
// This is mostly used by autogenerated code, and you should not need to use it manually.
func CheckAuthentication(authenticator Authenticator, ctx *Context) error {

	if authenticator == nil {
		return nil
	}

	ok, err := authenticator.IsAuthenticated(ctx)

	if err != nil {
		return err
	}

	if !ok {
		return elemental.NewError("Unauthorized", "You are not authorized to access this resource.", "bahamut", http.StatusUnauthorized)
	}

	return nil
}

// CheckAuthorization checks if the current context has been authorized if there is any authorizer registered.
//
// If it is not authorized it stops the normal processing execution flow, and will write the Unauthorized response to the given writer.
// If not Authorizer is set, then it will always return true.
//
// This is mostly used by autogenerated code, and you should not need to use it manually.
func CheckAuthorization(authorizer Authorizer, ctx *Context) error {

	if authorizer == nil {
		return nil
	}

	ok, err := authorizer.IsAuthorized(ctx)

	if err != nil {
		return elemental.NewError("Internal Server Error", err.Error(), "bahamut", http.StatusInternalServerError)
	}

	if !ok {
		return elemental.NewError("Forbidden", "You are not allowed to access this resource.", "bahamut", http.StatusForbidden)
	}

	return nil
}
