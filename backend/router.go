package backend

import (
	"fmt"
	"strings"

	"github.com/bitly/go-simplejson"
)

const (
	Read int = iota
	Write
)

// Request is used for Put/Get operations
type Request struct {
	// The following fields are used for routing
	Type int // Whether this is a Put/Read/Exists request
	MetaBlob bool // Whether the blob is a meta blob
	Host string
}

func (req *Request) String() string {
	return fmt.Sprintf("[request type=%v, meta=%v, hostname=%v]", req.Type, req.MetaBlob, req.Host)
}

type Router struct {
	Rules []*simplejson.Json

	Host string // Host used for rules checking

	Backends map[string]BlobHandler
}

// ResolveBackends construct the list of needed backend key
// by inspecting the rules
func (router *Router) ResolveBackends() []string {
	backends := []string{}
	for _, baseRule := range router.Rules {
		basicRule, err := baseRule.Array()
		_, basicMode := basicRule[0].(string)
		if err == nil && basicMode {
			// Basic rule handling [conf, backend]
			backends = append(backends, basicRule[1].(string))
		} else {
			backends = append(backends, baseRule.GetIndex(1).MustString())
		}
	}
	return backends
}

// TODO a way to set host

func (router *Router) Put(req *Request, hash string, data []byte) error {
	req.Type = Write
	key := router.Route(req)
	backend, exists := router.Backends[key]
	if !exists {
		panic(fmt.Errorf("backend %v is not registered", key))
	}
	return backend.Put(hash, data)
}

func (router *Router) Exists(req *Request, hash string) bool {
	req.Type = Read
	key := router.Route(req)
	backend, exists := router.Backends[key]
	if !exists {
		panic(fmt.Errorf("backend %v is not registered", key))
	}
	return backend.Exists(hash)
}

func (router *Router) Get(req *Request, hash string) (data []byte, err error) {
	req.Type = Read
	key := router.Route(req)
	backend, exists := router.Backends[key]
	if !exists {
		panic(fmt.Errorf("backend %v is not registered", key))
	}
	return backend.Get(hash)
}

func (router *Router) Enumerate(req *Request, res chan<- string) error {
	req.Type = Read
	key := router.Route(req)
	backend, exists := router.Backends[key]
	if !exists {
		panic(fmt.Errorf("backend %v is not registered", key))
	}
	return backend.Enumerate(res)
}

func (router *Router) Close() {
	for _, backend := range router.Backends {
		backend.Close()
	}
}

func (router *Router) Done() error {
	for _, backend := range router.Backends {
		if err := backend.Done(); err != nil {
			return err
		}
	}
	return nil
}

func NewRouterFromConfig(json *simplejson.Json) (*Router, error) {
	rules, err := json.Array()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config, body must be an array")
	}
	rconf := &Router{Rules: []*simplejson.Json{}, Backends: make(map[string]BlobHandler)}
	for i, _ := range rules {
		rconf.Rules = append(rconf.Rules, json.GetIndex(i))
	}
	return rconf, nil
}

// Route the request and return the backend key that match the request
func (router *Router) Route(req *Request) string {
	for _, baseRule := range router.Rules {
		basicRule, err := baseRule.Array()
		_, basicMode := basicRule[0].(string)
		if err == nil && basicMode {
			// Basic rule handling [conf, backend]
			backend := basicRule[1].(string)
			rule := basicRule[0].(string)
			if checkRule(rule, req) {
				return backend
			}
		} else {
			backend := baseRule.GetIndex(1).MustString()
			subRules, err := baseRule.GetIndex(0).StringArray()
			if err != nil {
				panic(fmt.Errorf("bad rule %v", baseRule.GetIndex(0)))
			}
			match := true
			for _, rule := range subRules {
				if !checkRule(rule, req) && match {
					match = false
				}
			}
			if match {
				return backend
			}
		}
	}
	return ""
}

// checkRule check if the rule match the given Request
func checkRule(rule string, req *Request) bool {
	switch {
	case rule == "if-meta":
		if req.MetaBlob {
			return true
		}
	case strings.HasPrefix(rule, "if-host-"):
		host := strings.Replace(rule, "if-host-", "", 1)
		if strings.ToLower(req.Host) == strings.ToLower(host) {
			return true
		}
	case rule == "default":
		return true
	default:
		panic(fmt.Errorf("failed to parse rule \"%v\"", rule))
	}
	return false
}