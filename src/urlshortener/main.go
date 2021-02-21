package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const DefaultPort = 8080
const DefaultNotFoundMessage = "Page %s is not found :("

type Config struct {
	NotFoundMessage string            `json:"not_found_message"`
	Port            int               `json:"server_port"`
	Paths           map[string]string `json:"paths"`
}

func ReadConfig(path string) (*Config, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Open file %s error: %s\n", path, err))
	}
	defer func() {
		err := jsonFile.Close()
		if err != nil {
			panic(fmt.Sprintf("Close config file %s error: %s", path, err))
		}
	}()

	rawConfig, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Read file %s error: %s\n", path, err))
	}

	config := Config{"", 0, map[string]string{}}
	err = json.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Parse json %s error: %s\n", string(rawConfig), err))
	}

	return &config, nil
}

func EnsurePortValid(config *Config) {
	if config.Port <= 0 {
		config.Port = DefaultPort
	}
}

func EnsureNotFoundMessageValid(config *Config) {
	if config.NotFoundMessage == "" {
		config.NotFoundMessage = DefaultNotFoundMessage
	}
}

func EnsureConfigValid(config *Config) {
	EnsurePortValid(config)
	EnsureNotFoundMessageValid(config)
}

func GetConfigPath() string {
	configPathPtr := flag.String("f", "data/url-mapping.json", "a file path")
	flag.Parse()
	return *configPathPtr
}

func RedirectHandler(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

type ProxyHttpHandler struct {
	notFoundMessage       string
	handlers              map[string]http.HandlerFunc
	availablePathsMessage string
}

func (h ProxyHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.handlers[r.URL.Path] == nil {
		h.notFoundHandler(w, r)
		return
	}

	h.handlers[r.URL.Path](w, r)
}

func (h ProxyHttpHandler) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	notFoundMessage := fmt.Sprintf(h.notFoundMessage, r.URL.Path)
	_, err := fmt.Fprintf(w, "%s\n\n%s", notFoundMessage, h.availablePathsMessage)
	if err != nil {
		fmt.Printf("Write http response error: %s", err)
		return
	}
}

func buildAvailablePathsMessage(config Config) string {
	message := "Available paths:\n"
	for short, long := range config.Paths {
		message = fmt.Sprintf("%s * %s -> %s\n", message, short, long)
	}

	return message
}

func PrepareHttpHandler(config Config) http.Handler {
	hh := ProxyHttpHandler{
		config.NotFoundMessage,
		map[string]http.HandlerFunc{},
		buildAvailablePathsMessage(config),
	}
	for short, long := range config.Paths {
		hh.handlers[short] = RedirectHandler(long)
	}
	return hh
}

func StartHttpServer(port int, hh http.Handler) error {
	fmt.Printf("Starting http server at http://localhost:%d ...\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), hh)
}

func main() {
	configPath := GetConfigPath()

	config, err := ReadConfig(configPath)
	if err != nil {
		fmt.Printf("Error parsing config json: %s", err)
		return
	}

	EnsureConfigValid(config)

	hh := PrepareHttpHandler(*config)

	err = StartHttpServer(config.Port, hh)
	if err != nil {
		fmt.Printf("Start Http Server error: %s", err)
		return
	}
}
