package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"ProjectGoLive/pkg/forms"
	"ProjectGoLive/pkg/models"

	"github.com/gorilla/mux"
	resty "gopkg.in/resty.v1"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	app.render(w, r, "home.page.tmpl", nil)
}

func (app *application) addPdtsvcForm(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Add Pdtsvc Form\n"))

	app.render(w, r, "addPdtsvc.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) editPdtsvcForm(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Edit Pdtsvc Form\n"))

	vars := mux.Vars(r)
	id := vars["id"]
	response, err := resty.R().Get(app.webSvcHost + "/api/v1/pdtsvcs/" + id)
	if err != nil {
		app.errorLog.Print("error getting data from the api service :: ", err)
		app.serverError(w, err)
		return
	}

	pdtsvc := models.Pdtsvc{}
	if err := json.NewDecoder(strings.NewReader(response.String())).Decode(&pdtsvc); err != nil {
		app.errorLog.Print("error occurred while decoding pdtsvc data :: ", err)
		app.serverError(w, err)
		return
	}

	app.render(w, r, "editPdtsvc.page.tmpl", &templateData{
		Pdtsvc: pdtsvc,
	})
}

func (app *application) addPdtsvc(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Add Pdtsvc\n"))

	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("code", "name", "description")
	form.MaxLength("code", 10)
	form.MaxLength("name", 100)

	pattern := `^[\D]+[\d]{3}$`
	myExp := regexp.MustCompile(pattern)
	form.MatchesPattern("code", myExp)

	if !form.Valid() {
		//fmt.Fprint(w, errors)
		app.render(w, r, "addPdtsvc.page.tmpl", &templateData{
			//FormErrors: errors,
			//FormData:   r.PostForm,
			Form: form,
		})
		return
	}

	pdtsvc := models.Pdtsvc{
		ID:          0,
		Code:        form.Get("code"),
		Name:        form.Get("name"),
		Description: form.Get("description"),
	}

	jsonMsg, err := json.Marshal(pdtsvc)
	if err != nil {
		app.errorLog.Print("ERROR occurred while encoding pdtsvc data :: ", err)
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("adding pdtsvc code :: %s with name"+
		" as :: %s and description as :: %s", pdtsvc.Code,
		pdtsvc.Name, pdtsvc.Description)

	response, err := resty.R().SetBody(jsonMsg).Post(app.webSvcHost + "/api/v1/pdtsvcs/create")
	if err != nil {
		app.errorLog.Print("error creating pdtsvc from the api service :: ", err)
		app.serverError(w, err)
		return
	}
	if err := json.NewDecoder(strings.NewReader(response.String())).Decode(&pdtsvc); err != nil {
		app.errorLog.Print("error occurred while decoding pdtsvc data :: ", err)
		app.serverError(w, err)
		return
	} else {
		//fmt.Fprintln(w, &pdtsvc)

		//app.session.Put(r, "flash", "Pdtsvc successfully created!")

		http.Redirect(w, r, "/pdtsvcs", http.StatusSeeOther)
	}
}

func (app *application) updatePdtsvc(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Update Pdtsvc\n"))

	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("code", "name", "description")
	form.MaxLength("code", 10)
	form.MaxLength("name", 100)

	pattern := `^[\D]+[\d]{3}$`
	myExp := regexp.MustCompile(pattern)
	form.MatchesPattern("code", myExp)

	if !form.Valid() {
		//fmt.Fprint(w, errors)
		app.render(w, r, "editPdtsvc.page.tmpl", &templateData{
			//FormErrors: errors,
			//FormData:   r.PostForm,
			Form: form,
		})
		return
	}

	intID, err := strconv.Atoi(form.Get("id"))
	if err != nil {
		app.errorLog.Print("ERROR occurred while decoding pdtsvc data :: ", err)
		app.serverError(w, err)
		return
	}

	pdtsvc := models.Pdtsvc{
		ID:          intID,
		Code:        form.Get("code"),
		Name:        form.Get("name"),
		Description: form.Get("description"),
	}

	jsonMsg, err := json.Marshal(pdtsvc)
	if err != nil {
		app.errorLog.Print("ERROR occurred while encoding pdtsvc data :: ", err)
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("updating pdtsvc id :: %d with"+
		" pdtsvc code :: %s with name"+
		" as :: %s and description as :: %s", pdtsvc.ID,
		pdtsvc.Code, pdtsvc.Name, pdtsvc.Description)

	response, err := resty.R().SetBody(jsonMsg).Put(app.webSvcHost + "/api/v1/pdtsvcs/update")
	if err != nil {
		app.errorLog.Print("error creating pdtsvc from the api service :: ", err)
		app.serverError(w, err)
		return
	}
	if err := json.NewDecoder(strings.NewReader(response.String())).Decode(&pdtsvc); err != nil {
		app.errorLog.Print("error occurred while decoding pdtsvc data :: ", err)
		app.serverError(w, err)
		return
	} else {
		//fmt.Fprintln(w, &pdtsvc)

		//app.session.Put(r, "flash", "Pdtsvc successfully updated!")

		http.Redirect(w, r, "/pdtsvcs", http.StatusSeeOther)
	}
}

func (app *application) deletePdtsvc(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Delete Pdtsvc\n"))

	vars := mux.Vars(r)
	id := vars["id"]
	response, err := resty.R().Get(app.webSvcHost + "/api/v1/pdtsvcs/" + id)
	if err != nil {
		app.errorLog.Print("error getting data from the api service :: ", err)
		app.serverError(w, err)
		return
	}

	pdtsvc := models.Pdtsvc{}
	if err := json.NewDecoder(strings.NewReader(response.String())).Decode(&pdtsvc); err != nil {
		app.errorLog.Print("error occurred while decoding pdtsvc data :: ", err)
		app.serverError(w, err)
		return
	}

	jsonMsg, err := json.Marshal(pdtsvc)
	if err != nil {
		app.errorLog.Print("ERROR occurred while encoding pdtsvc data :: ", err)
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("deleting pdtsvc id :: %d with"+
		" pdtsvc code :: %s with name"+
		" as :: %s and description as :: %s", pdtsvc.ID,
		pdtsvc.Code, pdtsvc.Name, pdtsvc.Description)

	_, err = resty.R().SetBody(jsonMsg).Delete(app.webSvcHost + "/api/v1/pdtsvcs/delete")
	if err != nil {
		app.errorLog.Print("error deleting pdtsvc from the api service :: ", err)
		app.serverError(w, err)
		return
	}

	//app.session.Put(r, "flash", "Pdtsvc successfully deleted!")

	http.Redirect(w, r, "/pdtsvcs", http.StatusSeeOther)
}

func (app *application) viewPdtsvc(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("View Pdtsvc\n"))

	vars := mux.Vars(r)
	id := vars["id"]
	response, err := resty.R().Get(app.webSvcHost + "/api/v1/pdtsvcs/" + id)
	if err != nil {
		app.errorLog.Print("error getting data from the api service :: ", err)
		app.serverError(w, err)
		return
	}
	//fmt.Fprintln(w, response.String())

	pdtsvc := models.Pdtsvc{}
	if err := json.NewDecoder(strings.NewReader(response.String())).Decode(&pdtsvc); err != nil {
		app.errorLog.Print("error occurred while decoding pdtsvc data :: ", err)
		app.serverError(w, err)
		return
	}

	app.render(w, r, "viewPdtsvc.page.tmpl", &templateData{
		Pdtsvc: pdtsvc,
	})

}

func (app *application) viewAllPdtsvcs(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("View All Pdtsvcs\n"))

	response, err := resty.R().Get(app.webSvcHost + "/api/v1/pdtsvcs")
	if err != nil {
		app.errorLog.Print("error getting data from the api service :: ", err)
		app.serverError(w, err)
		return
	}
	//fmt.Fprintln(w, response.String())

	pdtsvcs := []*models.Pdtsvc{}
	if err := json.NewDecoder(strings.NewReader(response.String())).Decode(&pdtsvcs); err != nil {
		app.errorLog.Print("error occurred while decoding pdtsvcs data :: ", err)
		app.serverError(w, err)
		return
	}

	app.render(w, r, "pdtsvcs.page.tmpl", &templateData{
		Pdtsvcs: pdtsvcs,
	})
}
