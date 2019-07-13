//go:binary-only-package

package lib

/*
	Some HTML presented at http://machine:port/debug/rpc
	Lists services, their methods, and some statistics, still rudimentary.
*/

import (
	"fmt"
	"net/http"
	"sort"
	"text/template"
)
