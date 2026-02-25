package mgp

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
)

const fileName = "goswag.go"

func GenerateSwagger(path string, routes []*Route, groups []*Group, defaultResponses []*ReturnType, defaultUseApiKeyAuth bool) {
	var (
		packagesToImport = make(map[string]bool)
		fullFileContent  = new(strings.Builder)
	)

	if path == "" {
		path = "."
	}

	routes, groups = addDefaultResponses(routes, groups, defaultResponses)
	routes, groups = addDefaultUseApiKeyAuth(routes, groups, defaultUseApiKeyAuth)

	if routes != nil {
		writeRoutes("", routes, fullFileContent, packagesToImport)
	}

	if groups != nil {
		writeGroup(groups, fullFileContent, packagesToImport)
	}

	f, err := os.Create(fmt.Sprintf("%s/%s", path, fileName))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writeFileContent(f, fullFileContent.String(), packagesToImport)
}

func addDefaultResponses(routes []*Route, groups []*Group, defaultResponses []*ReturnType) ([]*Route, []*Group) {
	if len(defaultResponses) == 0 {
		return routes, groups
	}

	for i := range routes {
		routes[i].Returns = append(routes[i].Returns, defaultResponses...)
	}

	for i := range groups {
		groups[i].Routes, groups[i].Groups = addDefaultResponses(groups[i].Routes, groups[i].Groups, defaultResponses)
	}

	return routes, groups
}

func addDefaultUseApiKeyAuth(routes []*Route, groups []*Group, defaultUseApiKeyAuth bool) ([]*Route, []*Group) {
	if !defaultUseApiKeyAuth {
		return routes, groups
	}

	for i := range routes {
		routes[i].UseApiKeyAuth = defaultUseApiKeyAuth
	}

	for i := range groups {
		groups[i].Routes, groups[i].Groups = addDefaultUseApiKeyAuth(groups[i].Routes, groups[i].Groups, defaultUseApiKeyAuth)
	}

	return routes, groups
}

func writeFileContent(file io.Writer, content string, packagesToImport map[string]bool) {
	_, err := fmt.Fprintf(file, "package main\n\n")
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if len(packagesToImport) > 0 {
		_, err = fmt.Fprintf(file, "import (\n")
		if err != nil {
			log.Fatalf("%+v", err)
		}

		for pkg := range packagesToImport {
			_, err = fmt.Fprintf(file, "\t_ \"%s\"\n", pkg)
			if err != nil {
				log.Fatalf("%+v", err)
			}
		}

		_, err = fmt.Fprintf(file, ")\n\n")
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}

	_, err = fmt.Fprintf(file, "%s", content)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func writeRoutes(groupName string, routes []*Route, s *strings.Builder, packagesToImport map[string]bool) {
	for i := range routes {
		addLineIfNotEmpty(s, routes[i].Summary, "// @Summary %s\n")
		addTextIfNotEmptyOrDefault(s, routes[i].Summary, "// @Description %s\n", routes[i].Description)

		if len(routes[i].Tags) > 0 {
			s.WriteString(fmt.Sprintf("// @Tags %s\n", strings.Join(routes[i].Tags, ",")))
		} else if groupName != "" {
			s.WriteString(fmt.Sprintf("// @Tags %s\n", groupName))
		}

		if routes[i].Method == http.MethodPost || routes[i].Method == http.MethodPut {
			addTextIfNotEmptyOrDefault(s, "json", "// @Accept %s\n", routes[i].Accepts...)
		}

		if routes[i].Returns != nil {
			addTextIfNotEmptyOrDefault(s, "json", "// @Produce %s\n", routes[i].Produces...)
		}

		if routes[i].QueryStruct != nil {
			s.WriteString(fmt.Sprintf("// @Param request query %s true \"Query\"\n", getStructAndPackageName(routes[i].QueryStruct)))
		}

		if routes[i].BodyStruct != nil {
			s.WriteString(fmt.Sprintf("// @Param request body %s true \"Request\"\n", getStructAndPackageName(routes[i].BodyStruct)))
		}

		if routes[i].PathStruct != nil {
			s.WriteString(fmt.Sprintf("// @Param request path %s true \"Path\"\n", getStructAndPackageName(routes[i].PathStruct)))
		}

		if routes[i].HeaderStruct != nil {
			s.WriteString(fmt.Sprintf("// @Param request header %s true \"Header\"\n", getStructAndPackageName(routes[i].HeaderStruct)))
		}

		if routes[i].Returns != nil {
			writeReturns(routes[i].Returns, s, packagesToImport)
		}

		if routes[i].UseApiKeyAuth {
			s.WriteString("// @Security ApiKeyAuth\n")
		}

		if routes[i].Path != "" {
			sampleRegexp := regexp.MustCompile(`:(\w+)`)
			s.WriteString(fmt.Sprintf("// @Router %s [%s]\n", sampleRegexp.ReplaceAllString(routes[i].Path, "{$1}"), strings.ToLower(routes[i].Method)))
		}

		if routes[i].FuncName != "" {
			s.WriteString(fmt.Sprintf("func %s() {} //nolint:unused \n", routes[i].FuncName))
		}

		s.WriteString("\n")
	}
}

func writeReturns(returns []*ReturnType, s *strings.Builder, packagesToImport map[string]bool) {
	for i := range returns {
		if returns[i].StatusCode == 0 {
			continue
		}

		respType := "@Success"
		firstDigit := returns[i].StatusCode / 100

		if firstDigit != http.StatusOK/100 { // <> 2xx
			respType = "@Failure"
		}

		if returns[i].Body == nil {
			s.WriteString(fmt.Sprintf("// %s %d\n", respType, returns[i].StatusCode))
			continue
		}

		isGeneric := writeIfIsGenericType(s, returns[i], respType)

		if !isGeneric {
			s.WriteString(fmt.Sprintf("// %s %d {object} %s", respType, returns[i].StatusCode, getStructAndPackageName(returns[i].Body)))
		}

		addPackageToImport(returns[i], packagesToImport)

		s.WriteString("\n")
	}
}

func writeGroup(groups []*Group, s *strings.Builder, packagesToImport map[string]bool) {
	for i := range groups {
		writeRoutes(groups[i].GroupName, groups[i].Routes, s, packagesToImport)

		if groups[i].Groups != nil {
			writeGroup(groups[i].Groups, s, packagesToImport)
		}
	}
}

func addPackageToImport(data *ReturnType, packagesToImport map[string]bool) {
	if data.Body == nil {
		return
	}
	t := reflect.TypeOf(data.Body)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.PkgPath() != "" {
		packagesToImport[t.PkgPath()] = true
	}
}

func writeIfIsGenericType(s *strings.Builder, data *ReturnType, respType string) (isGeneric bool) {
	bodyName := getStructAndPackageName(data.Body)
	isGeneric = bodyName[len(bodyName)-1:] == "]"
	if !isGeneric {
		return
	}

	bodyName = regexp.MustCompile(`\[(\w+/)*(\w+)\.`).ReplaceAllString(bodyName, "[$2.")
	bodyName = regexp.MustCompile(`](\w+/)*(\w+)\.`).ReplaceAllString(bodyName, "]$2.")

	s.WriteString(fmt.Sprintf("// %s %d {object} %s", respType, data.StatusCode, bodyName))
	return isGeneric
}

func getStructAndPackageName(body any) string {
	isPointer := reflect.TypeOf(body).Kind() == reflect.Ptr
	if isPointer {
		body = reflect.ValueOf(body).Elem().Interface()
	}

	return reflect.TypeOf(body).String()
}

func addTextIfNotEmptyOrDefault(s *strings.Builder, defaultText, format string, text ...string) {
	if text != nil {
		if len(text) >= 1 && strings.TrimSpace(text[0]) != "" {
			s.WriteString(fmt.Sprintf(format, strings.Join(text, ",")))
			return
		}
	}

	if defaultText != "" {
		s.WriteString(fmt.Sprintf(format, defaultText))
	}
}

func addLineIfNotEmpty(s *strings.Builder, data, format string) {
	if data != "" {
		s.WriteString(fmt.Sprintf(format, data))
	}
}
