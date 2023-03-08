package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/nippati/moff/pkg/ansible"
	"github.com/nippati/moff/pkg/scan"
	"github.com/nippati/moff/pkg/ui"
)

type Vulnerability struct {
	ID string `json:"id"`
}

func main() {
	// Read the JSON file of vuls scan results
	file, err := os.Open("vuls-results.json")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	var results scan.Results
	err = json.NewDecoder(file).Decode(&results)
	if err != nil {
		log.Fatalf("error decoding JSON: %v", err)
	}

	// Convert []scan.Vulnerability to []ui.Vulnerability
	var vulns []ui.Vulnerability
	for _, r := range results {
		v := ui.Vulnerability{
			ID: r.VulnerabilityID,
		}
		vulns = append(vulns, v)
	}

	// Define the HTML template for the vulnerabilities page
	tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Vulnerabilities</title>
		</head>
		<body>
			<h1>Select Vulnerabilities</h1>
			<form action="/" method="POST">
				{{range .}}
				<input type="checkbox" name="vulns" value="{{.ID}}">{{.ID}}<br>
				{{end}}
				<input type="submit" value="Generate Ansible Playbook">
			</form>
		</body>
		</html>
	`
	// Parse the template
	t, err := template.New("vulns").Parse(tmpl)
	if err != nil {
		log.Fatalf("error parsing template: %v", err)
	}

	// Create a HTTP server and listen for requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// User submitted the form, process the selected vulnerabilities
			r.ParseForm()
			selectedIDs := r.Form["vulns"]

			var selected []ui.Vulnerability
			for _, v := range vulns {
				for _, id := range selectedIDs {
					if v.ID == id {
						selected = append(selected, v)
					}
				}
			}

			if len(selected) == 0 {
				fmt.Fprintln(w, "No vulnerabilities selected!")
				return
			}

			playbook := ansible.GeneratePlaybook(selected)

			// Write the playbook to a file
			file, err := os.Create("playbook.yml")
			if err != nil {
				fmt.Fprintln(w, "Error creating playbook file:", err)
				return
			}
			defer file.Close()

			_, err = file.WriteString(playbook)
			if err != nil {
				fmt.Fprintln(w, "Error writing playbook file:", err)
				return
			}

			// Return a success message to the user
			w.Write([]byte("Ansible playbook successfully generated!"))
			// Return a success message to the user
			w.Write([]byte("Ansible playbook successfully generated!"))

			// Execute the Ansible playbook
			err = ansible.RunPlaybook("playbook.yml")
			if err != nil {
				fmt.Println("Error executing playbook:", err)
				return
			}

			// Return a success message to the user
			w.Write([]byte("Ansible playbook successfully executed!"))
		} else {
			// User requested the vulnerabilities page
			err := t.Execute(w, vulns)
			if err != nil {
				fmt.Fprintln(w, "Error executing template:", err)
				return
			}
		}

	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// GeneratePlaybook generates an Ansible playbook based on the selected vulnerabilities
func GeneratePlaybook(vulns []ui.Vulnerability) string {
	// TODO: implement playbook generation logic
	return ""
}

// RunPlaybook runs the specified Ansible playbook
func RunPlaybook(playbook string) error {
	// TODO: implement playbook execution logic
	return nil
}
