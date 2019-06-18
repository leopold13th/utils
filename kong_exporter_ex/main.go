package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getJSONForURL(urlPath string) map[string]interface{} {
	url := "http://localhost:8001" + urlPath
	// fmt.Println("URL:>", url)

	// var jsonStr = []byte(`{"title": "Buy cheese and bread for breakfast."}`)
	var jsonStr = []byte(`{}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	jsonData := []byte(string(body))
	var v interface{}
	json.Unmarshal(jsonData, &v)
	data := v.(map[string]interface{})

	return data
}

func getRoot() {
	data := getJSONForURL("/")

	data = data["plugins"].(map[string]interface{})
	var pluginsEnabled = data["enabled_in_cluster"].([]interface{})

	for _, v := range pluginsEnabled {
		// fmt.Println(v)
		pluginsEnabledGauge.WithLabelValues(v.(string)).Set(1)
	}
}

func getServices() {
	data := getJSONForURL("/services")

	services := data["data"].([]interface{})
	for _, service := range services {
		service := service.(map[string]interface{})
		// fmt.Println(service["name"])
		servicesGauge.WithLabelValues(service["name"].(string)).Set(1)
	}
}

// func getRoutes() {
// 	data := getJSONForURL("/routes")

// 	routes := data["data"].([]interface{})
// 	for _, route := range routes {
// 		route := route.(map[string]interface{})
// 		fmt.Println(route["paths"].([]interface{}))
// 		routesGauge.WithLabelValues(route["paths"].([]interface{})).Set(1)
// 	}
// }

func getUpstreams() {
	data := getJSONForURL("/upstreams")

	upstreams := data["data"].([]interface{})
	for _, upstream := range upstreams {
		upstream := upstream.(map[string]interface{})

		data := getJSONForURL("/upstreams/" + upstream["name"].(string) + "/targets")
		targets := data["data"].([]interface{})
		for _, target := range targets {
			target := target.(map[string]interface{})
			// fmt.Println(upstream["name"], "--> ", target["target"])
			upstreamsGauge.WithLabelValues(upstream["name"].(string), target["target"].(string)).Set(1)
		}
	}
}

func recordMetrics() {
	go func() {
		for {
			getRoot()
			getServices()
			// getRoutes()
			getUpstreams()
			// opsProcessed.Inc()
			time.Sleep(30 * time.Second)
		}
	}()
}

var (
	// opsProcessed = promauto.NewCounter(
	// 	prometheus.CounterOpts{
	// 		Name: "myapp_processed_ops_total",
	// 		Help: "The total number of processed events",
	// 	})

	pluginsEnabledGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kong_plugins_enabled",
			Help: "Enabled plugins",
		},
		[]string{"plugin"})

	servicesGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kong_services",
			Help: "Services configured in Kong",
		},
		[]string{"service"})

	routesGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kong_routes",
			Help: "Routes",
		},
		[]string{"route"})

	upstreamsGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kong_upstreams",
			Help: "Upstreams attached to services",
		},
		[]string{"upstream", "target"})
)

func main() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9080", nil)
}
