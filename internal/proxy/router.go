package proxy

import (
	"strings"

	"github.com/david-tobi-peter/Prism/internal/config"
)

type Router struct {
	cfg *config.GatewayConfig
}

func NewRouter(cfg *config.GatewayConfig) *Router {
	return &Router{cfg: cfg}
}

func (r *Router) Resolve(inboundPath string) (string, string, bool) {
	for _, route := range r.cfg.Routes {
		if strings.HasPrefix(inboundPath, route.Path) {
			workingPath := inboundPath

			if route.StripPrefix {
				workingPath = strings.TrimPrefix(workingPath, route.Path)
				if !strings.HasPrefix(workingPath, "/") {
					workingPath = "/" + workingPath
				}
			}

			for _, rule := range route.Rewrites {
				if rule.Pattern.MatchString(workingPath) {
					workingPath = rule.Pattern.ReplaceAllString(workingPath, rule.Replace)
					break
				}
			}

			return route.Backend, workingPath, true
		}
	}

	return "", "", false
}
