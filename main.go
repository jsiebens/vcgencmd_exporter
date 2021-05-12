package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"vcgencmd_exporter/pkg"
)

func main() {

	config := pkg.LoadConfig()

	throttled := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rpi_vcgencmd_get_throttled",
	}, []string{"bit", "description"})

	checks := map[int64]string{
		0:  "Under-voltage detected",
		1:  "Arm frequency capped",
		2:  "Currently throttled",
		3:  "Soft temperature limit active",
		16: "Under-voltage has occurred",
		17: "Arm frequency capped has occurred",
		18: "Throttling has occurred",
		19: "Soft temperature limit has occurred",
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	}()

	for {
		select {
		case <-signalChan:
			return
		default:
			cmd := exec.Command(config.VcgencmdBinary, "get_throttled")
			stdout, err := cmd.Output()

			if err != nil {
				fmt.Println(err)
			} else {
				v := strings.TrimPrefix(string(stdout), "throttled=")
				if s, err := strconv.ParseInt(v, 0, 32); err == nil {
					for key, value := range checks {
						throttled.With(prometheus.Labels{
							"bit":         fmt.Sprintf("%v", key),
							"description": value,
						}).Set(float64(checkBit(s, key+1)))
					}
				}
			}

			time.Sleep(time.Duration(config.Delay) * time.Second)
		}
	}
}

func checkBit(n, k int64) int64 {
	return (n >> (k - 1)) & 1
}
