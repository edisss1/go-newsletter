package server

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/go-mail/mail/v2"
	"google.golang.org/api/option"
)

func FormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	serviceAccountKeyJSON := os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY")
	if serviceAccountKeyJSON == "" {
		http.Error(w, "FIREBASE_SERVICE_ACCOUNT_KEY environment variable is not set", http.StatusInternalServerError)
		return
	}

	opt := option.WithCredentialsJSON([]byte(serviceAccountKeyJSON))

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Firestore client: %v\n", err)
	}

	defer client.Close()
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusMethodNotAllowed)
	}

	tmpl := template.Must(template.ParseFiles("templates/submit.html"))
	tmpl.Execute(w, nil)

	subject := r.FormValue("subject")
	message := r.FormValue("message")

	recipients, err := FetchDocuments(client, ctx)
	if err != nil {
		log.Fatalf("Error fetching documents")
	}

	if err := sendEmail(subject, message, recipients); err != nil {
		log.Printf("Error sending email: %v", err)
		http.Error(w, "Error sending email", http.StatusInternalServerError)
		return
	}

	log.Printf("Subject: %v, Message: %v", subject, message)
}

func sendEmail(subject, message string, recipients []string) error {
	m := mail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_ADDRESS"))
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)

	d := mail.NewDialer("smtp.gmail.com", 465, os.Getenv("EMAIL_ADDRESS"), os.Getenv("EMAIL_PASSWORD"))
	d.SSL = true
	for _, recipient := range recipients {
		m.SetHeader("To", recipient)
		if err := d.DialAndSend(m); err != nil {
			return err
		}
		log.Printf("Message sent to %s", recipient)
	}
	return nil
}
