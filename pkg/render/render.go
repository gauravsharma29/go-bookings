package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/gauravsharma29/go-bookings/pkg/config"
	"github.com/gauravsharma29/go-bookings/pkg/models"
)

var app *config.AppConfig

// sets the config for the template package
func NewTemplate(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

func RenderTemplate(w http.ResponseWriter, html string, td *models.TemplateData) {
	// create a template cache
	var tc map[string]*template.Template
	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache t= template
	t, ok := tc[html]
	if !ok {
		log.Fatal("Could not get the template from template cache")
	}

	// extra steps
	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing templates to browser")
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named with *.page.html from ./templates
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	// range through all the files ending with *.page.html
	for _, page := range pages {
		fileName := filepath.Base(page) // page = home.page.html
		ts, err := template.New(fileName).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// get all layouts
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			// parse layout file and add it to the existing ts to make it one to render
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[fileName] = ts
	}

	return myCache, nil
}

// -------------------

// var tc = make(map[string]*template.Template)

// func RenderTemplate(w http.ResponseWriter, t string) {
// 	var html *template.Template
// 	var err error

// 	// check to see if we already have the template in our cache
// 	if _, inMap := tc[t]; !inMap {
// 		// need to create a template
// 		log.Println("Creating template and adding to cache")
// 		err = createTemplateCache(t)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 	} else {
// 		// we have the template in the cache
// 		log.Println("using cached template")
// 	}

// 	html = tc[t]

// 	err = html.Execute(w, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

// func createTemplateCache(t string) error {
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", t),
// 		"./templates/base.layout.html",
// 	}

// 	// parse the template
// 	html, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err
// 	}

// 	// add template to cache
// 	tc[t] = html

// 	return nil
// }

// First version of renderingtemplate
// func RenderTemplateTest(w http.ResponseWriter, html string) {
// 	parsedTemplate, _ := template.ParseFiles("./templates/"+html, "./templates/base.layout.html")
// 	err := parsedTemplate.Execute(w, nil)
// 	if err != nil {
// 		fmt.Println("error parsing templates:", err)
// 		return
// 	}
// }
