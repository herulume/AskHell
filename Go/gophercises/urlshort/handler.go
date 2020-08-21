package urlshort

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

const fallbackResponse = "fallback"

func fallback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, fallbackResponse)
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathURLs []pathURL
	err := yaml.Unmarshal(yml, &pathURLs)
	if err != nil {
		return nil, err
	}

	pathToUrls := buildMap(pathURLs)
	return MapHandler(pathToUrls, fallback), nil
}

type pathURL struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func buildMap(pathURLs []pathURL) map[string]string {
	pathToUrls := make(map[string]string)

	for _, p := range pathURLs {
		pathToUrls[p.Path] = p.Url
	}

	return pathToUrls
}
