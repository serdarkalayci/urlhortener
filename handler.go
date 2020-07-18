package urlshortener

import (
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
		newURL := pathsToUrls[r.URL.EscapedPath()]
		if newURL != "" {
			http.Redirect(w, r, newURL, http.StatusSeeOther)
		}
		fallback.ServeHTTP(w, r)
	}
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
	// First unmarshal the yaml into an array of pathMap struct
	var pathMaps []pathMap
	err := yaml.Unmarshal(yml, &pathMaps)
	if err != nil {
		return nil, err
	}
	// Convert []pathMap to map[string]string so it can be used by MapHandler
	mapMap := make(map[string]string)
	for _, pm := range pathMaps {
		mapMap[pm.Path] = pm.URL
	}
	// Return MapHandler to serve this request, there's no need to fallback
	return MapHandler(mapMap, fallback), nil
}

type pathMap struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
