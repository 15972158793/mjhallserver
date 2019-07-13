//go:binary-only-package

package lib

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sync"
)
