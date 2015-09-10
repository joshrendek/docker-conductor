package healthcheck

import (
	log "gopkg.in/inconshreveable/log15.v2"
	"net/url"
	"os/exec"
	"strings"
)

type HealthCheck struct {
	Script string
	Log    log.Logger
	Host   string
}

func New(l log.Logger, script string, host string) *HealthCheck {
	u, _ := url.Parse(host)
	ip_host := strings.Split(u.Host, ":")[0]
	return &HealthCheck{Log: l, Script: script, Host: ip_host}
}

func (h *HealthCheck) Check() bool {
	script := strings.Split(h.Script, " ")
	for i, p := range script {
		if p == "$HOST" {
			script[i] = h.Host
		}
	}
	cmd := exec.Command(script[0], script...)
	cmd.Env = []string{"HOST=" + h.Host}
	h.Log.Info("[ ] Running health check: " + strings.Join(script, " "))
	out, err := cmd.CombinedOutput()
	if err != nil {
		h.Log.Error("[F] Health check failed")
		h.Log.Error(err.Error())
		h.Log.Error(string(out))
		return false
	} else {
		h.Log.Info("[P] Health check passed")
		return true
	}
}
