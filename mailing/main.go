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

func (t *TemplateReg) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// to := []string{
	// 	"20bce091@nith.ac.in",
	// }
	// message := []byte("This is a email body")
	// SendMail(to, message)
	e := echo.New()

	e.Renderer = &TemplateReg{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
  e.GET("/", Home )
	e.Logger.Fatal(e.Start(":8000"))
}


func Home(c echo.Context) error {
  return c.Render(http.StatusOK, "home.html", map[string]interface{}{
    "name": "Sumit",
    "msg": "Hello, lappy!",
  })
}

func SendMail(to []string, message []byte) {
	from := "sd08012003@gmail.com"
	pswd := "bjqxxtatmzeomwwe"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, pswd, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
	}
}
