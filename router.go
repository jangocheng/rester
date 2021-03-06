// Copyright 2020 HenryLee. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rester

import (
	"fmt"
	"reflect"

	"github.com/buaazp/fasthttprouter"
	"github.com/henrylee2cn/ameda"
)

type Router struct {
	router          fasthttprouter.Router
	controllers     map[string]Controller // {relativePath:Controller}
	controllerNames map[string]string     // {controllerName:relativePath}
}

// Control routes controller.
// NOTE:
// The same routing controller can be registered repeatedly, but only for the first time;
// If the controller of the same route registered twice is different, panic
func (r *Router) Control(path string, controller Controller) {
	if r.controllers == nil {
		r.controllers = make(map[string]Controller)
	}
	if r.controllerNames == nil {
		r.controllerNames = make(map[string]string)
	}
	ctl, ok := r.controllers[path]
	if ok && reflect.TypeOf(ctl) == reflect.TypeOf(controller) {
		return
	}
	handlerMap := MustToHandlers(controller)
	controllerName := getControllerName(controller)
	for _, httpMethod := range httpMethodList {
		handler := handlerMap[httpMethod]
		if handler != nil {
			r.router.Handle(httpMethod, path, handler)
			r.controllerNames[controllerName] = path
			r.println(httpMethod, path, controllerName)
		}
	}
	r.controllers[path] = controller
}

// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
//     router.ServeFiles("/src/*filepath", "/var/www")
func (r *Router) ServeFiles(path string, rootPath string) {
	r.router.ServeFiles(path, rootPath)
	r.println("GET", path, "fasthttp.FSHandler")
}

// Path returns router path of the controller
// NOTE:
//  Must be called after routing
func (r *Router) Path(controller Controller) string {
	return r.controllerNames[getControllerName(controller)]
}

func (r *Router) println(httpMethod, path, controllerName string) {
	fmt.Printf(
		"[RESTER] %-7s %-25s --> %s\n",
		httpMethod, path, controllerName,
	)
}

func getControllerName(controller Controller) string {
	t := ameda.DereferenceValue(reflect.ValueOf(controller)).Type()
	return t.PkgPath() + "." + t.Name()
}
