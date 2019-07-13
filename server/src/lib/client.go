//go:binary-only-package

package lib

import (
	"bufio"
	"encoding/gob"
	"errors"
	"io"
	"net"
	"net/http"
	"reflect"
	"sync"
)
