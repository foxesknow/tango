package xmlmacro

import (
	"strings"
)

type nameScope struct {
	previous  *nameScope
	variables map[string]any
}

// Create a new scope
func newNameScope() *nameScope {
	return &nameScope{
		previous:  nil,
		variables: make(map[string]any),
	}
}

// Creates a new scope, with the current scope becoming the previous scoope
func (scope *nameScope) newChildScope() *nameScope {
	return &nameScope{
		previous:  scope,
		variables: make(map[string]any),
	}
}

// Get the named value from this scope or any previous scopes
func (scope *nameScope) getValue(name string) (any, bool) {
	normalizedName := scope.normalizeName(name)

	activeScope := scope

	for activeScope != nil {
		if value, found := activeScope.variables[normalizedName]; found {
			return value, true
		}
		activeScope = activeScope.previous
	}

	return nil, false
}

// Declares a new variable in the current scope if it is not already declared
func (scope *nameScope) declare(name string, value any) bool {
	normalizedName := scope.normalizeName(name)

	if _, found := scope.variables[normalizedName]; found {
		return false
	}

	scope.variables[normalizedName] = value
	return true
}

// Checks to see if a variable is declared in this scope only
func (scope *nameScope) isDeclaredInActiveScope(name string) bool {
	normalizedName := scope.normalizeName(name)
	_, found := scope.variables[normalizedName]

	return found
}

func (scope *nameScope) normalizeName(name string) string {
	return strings.ToLower(name)
}
