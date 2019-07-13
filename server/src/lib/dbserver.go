//go:binary-only-package

package lib

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"rjmgr"
	"strings"
)
