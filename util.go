package mgp

import (
	"path"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

func getFullPath(groupName, relativePath string) string {
	if groupName == "" {
		return relativePath
	}

	fullPath := path.Join(groupName, relativePath)

	if strings.HasSuffix(relativePath, "/") {
		fullPath += "/"
	}

	return fullPath
}

func getFuncName(handlers ...func(c *Context)) string {
	lastHandler := handlers[len(handlers)-1]

	fullFuncName := runtime.FuncForPC(reflect.ValueOf(lastHandler).Pointer()).Name()
	funcNameSplit := strings.Split(fullFuncName, ".")
	funcName := funcNameSplit[len(funcNameSplit)-1]
	funcName = strings.TrimSuffix(funcName, "-fm")

	return funcName
}

func toGoSwagGroup(from []*RouterGroup) []*Group {
	var groups []*Group
	for i := range from {
		sg := &Group{
			GroupName: from[i].groupName,
			Routes:    from[i].routes,
		}
		if len(from[i].groups) > 0 {
			sg.Groups = toGoSwagGroup(from[i].groups)
		}
		groups = append(groups, sg)
	}

	return groups
}

func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
