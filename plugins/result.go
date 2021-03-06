package plugins

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/rbeuque74/nsca"
)

// StatusEnum corresponds to all status for a Consumer message
type StatusEnum int16

const (
	// STATE_OK represents a healthy service
	STATE_OK StatusEnum = nsca.STATE_OK
	// STATE_WARNING represents a service that requires attention
	STATE_WARNING StatusEnum = nsca.STATE_WARNING
	// STATE_CRITICAL represents a service that needs immediate fix
	STATE_CRITICAL StatusEnum = nsca.STATE_CRITICAL
	// STATE_UNKNOWN represents a service in an unsure status
	STATE_UNKNOWN StatusEnum = nsca.STATE_UNKNOWN

	errNilTemplate    = "unable to apply nil jagozzi template"
	errFailedTemplate = "unable to apply jagozzi template %q: %s"
)

// Result is the structure that represents a checker result
type Result struct {
	// Status indicates if check was successful or not
	Status StatusEnum
	// Message is the additional message that can be given to Consumer server
	Message string
	// Checker is the checker that returns this result
	Checker Checker
}

// RenderError allow personalised rendering if checker contains a template
func RenderError(tmpl *template.Template, model interface{}) string {
	if tmpl == nil {
		return errNilTemplate
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, model); err != nil {
		return fmt.Sprintf(errFailedTemplate, tmpl.Name(), err)
	}
	return buf.String()
}

// ResultFromError generates a critical result from a checker and an error
func ResultFromError(checker Checker, err error, prefix string) Result {
	msg := err.Error()
	if prefix != "" {
		msg = prefix + ": " + msg
	}

	return Result{
		Status:  STATE_CRITICAL,
		Message: msg,
		Checker: checker,
	}
}
