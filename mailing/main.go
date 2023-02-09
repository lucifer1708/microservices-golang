package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/smtp"

	"github.com/labstack/echo/v4"
)

type TemplateReg struct {
	templates *template.Template
}

type Formdata struct {
	To      string `json:"to" form:"to" query:"to"`
	Message string `json:"message" form:"message" query:"message"`
}

func main() {
	e := echo.New()
	e.Renderer = &TemplateReg{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.GET("/", Home)
	e.POST("/send", Form)
	e.Logger.Fatal(e.Start(":8000"))
}

// Template Renderer
func (t *TemplateReg) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Home page function
func Home(c echo.Context) error {
	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"msg": "fill this form",
	})
}

// function with post request to send an email with data filled into the form.
func Form(c echo.Context) error {
	formdata := new(Formdata)
	if err := c.Bind(formdata); err != nil {
		fmt.Println(err)
	}
	tomail := formdata.To
	to := []string{tomail}
	msg := formdata.Message
	SendMail(to, []byte(msg))
	return c.Render(http.StatusOK, "success.html", map[string]interface{}{
		"to":  to,
		"msg": msg,
	})

}

// Send Mail Function used to send emails
func SendMail(to []string, message []byte) {
	from := "example@gmail.com"
	pswd := "enteryourapppassword"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, pswd, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
	}
}
