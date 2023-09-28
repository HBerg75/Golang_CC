package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"time"
)

type Repository struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	CloneURL    string    `json:"clone_url"`
}

func main() {
	githubOrg := os.Getenv("GITHUB_ORG")
	githubUser := os.Getenv("GITHUB_USER")
	if githubOrg == "" && githubUser == "" {
		fmt.Println("Veuillez définir GITHUB_ORG ou GITHUB_USER")
		return
	}

	var url string
	if githubUser != "" {
		url = fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated", githubUser)
	} else {
		url = fmt.Sprintf("https://api.github.com/orgs/%s/repos?sort=updated", githubOrg)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var repos []Repository
	err = json.Unmarshal(body, &repos)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Trier les dépôts par date de dernière modification
	sort.SliceStable(repos, func(i, j int) bool {
		return repos[i].UpdatedAt.After(repos[j].UpdatedAt)
	})

	// Créer un fichier CSV
	file, err := os.Create("repositories.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Écrire l'en-tête du CSV
	err = writer.Write([]string{"Name", "Description", "Updated At", "Clone URL"})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Écrire les données des dépôts dans le CSV
	for _, repo := range repos {
		err := writer.Write([]string{repo.Name, repo.Description, repo.UpdatedAt.String(), repo.CloneURL})
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err = os.Mkdir("clones", os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}


		// ...
		err = os.Mkdir("clones", os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	
		token := os.Getenv("GITHUB_TOKEN")  // Assurez-vous de définir cette variable d'environnement
	
		for _, repo := range repos {
			cmd := exec.Command("git", "clone", fmt.Sprintf("https://%s:x-oauth-basic@github.com/user/%s.git", token, repo.Name))
			cmd.Dir = "clones"  // Définir le répertoire dans lequel la commande sera exécutée
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Erreur lors du clonage du dépôt %s: %v\n", repo.Name, err)
			} else {
				fmt.Printf("Dépôt %s cloné avec succès\n", repo.Name)
			}
		}
		// ...
	

	
}
