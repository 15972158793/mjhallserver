//go:binary-only-package

package lib

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	//"log"
	"net"
	"net/http"
	"reflect"
	"rjmgr"
	sysdebug "runtime/debug"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)
