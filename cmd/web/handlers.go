package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"vtorosyan.learning/internal/models"
	"vtorosyan.learning/internal/validator"
)

const (
	ErrTitleInvalid   = "title can not be blank"
	ErrTitleTooLong   = "title should be less than 100 characters"
	ErrContentInvalid = "content can not be blank"
	ErrExpiresInvalid = "expire field must equal 1, 7 or 365"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, r, err)
		}
	}

	tData := app.newTemplateData(r)
	tData.Snippets = snippets
	app.render(w, r, http.StatusOK, "home.tmpl.html", tData)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusNotFound)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	tData := app.newTemplateData(r)
	tData.Snippet = snippet
	app.render(w, r, http.StatusOK, "view.tmpl.html", tData)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	tData := app.newTemplateData(r)
	tData.Form = snippetCreateForm{Expires: 365}
	app.render(w, r, http.StatusOK, "create.tmpl.html", tData)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var snippetForm snippetCreateForm

	err = app.decodePostForm(r, &snippetForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	snippetForm.CheckField(validator.NotBlank(snippetForm.Title), "title", ErrTitleInvalid)
	snippetForm.CheckField(validator.MaxChars(snippetForm.Title, 100), "title", ErrTitleTooLong)
	snippetForm.CheckField(validator.NotBlank(snippetForm.Content), "content", ErrContentInvalid)
	snippetForm.CheckField(validator.PermittedValue(snippetForm.Expires, 1, 7, 365), "expires", ErrExpiresInvalid)

	if !snippetForm.Valid() {
		data := app.newTemplateData(r)
		data.Form = snippetForm
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(snippetForm.Title, snippetForm.Content, snippetForm.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
