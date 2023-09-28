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
	"strings"
	"archive/zip"
	"path/filepath"
)

type Repository struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	CloneURL    string    `json:"clone_url"`
}


func handleArchive(w http.ResponseWriter, r *http.Request) {
	// Vérifiez la méthode HTTP
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	zipFile, err := os.Create("repositories.zip")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Ajouter des dossiers et des fichiers à l'archive ZIP
	err = filepath.Walk("clones", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = path
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fileContent, err := os.ReadFile(path)  
		if err != nil {
			return err
		}
		_, err = writer.Write(fileContent)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Archive créée avec succès"))
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Erreur : %s\n", body)
		return
	}	

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
	
		token := os.Getenv("GITHUB_TOKEN") 
	
		for _, repo := range repos {
			owner := githubUser
			if owner == "" {
				owner = githubOrg
			}
			cmd := exec.Command("git", "clone", fmt.Sprintf("https://%s:x-oauth-basic@github.com/%s/%s.git", token, owner, repo.Name))
			cmd.Dir = "clones"
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Erreur lors du clonage du dépôt %s: %v\n", repo.Name, err)
				continue  // Continue avec le prochain repo si le clonage échoue 
			} 
			fmt.Printf("Dépôt %s cloné avec succès\n", repo.Name)

			repoDir := fmt.Sprintf("clones/%s", repo.Name)
			cmd = exec.Command("git", "-C", repoDir, "fetch")
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Erreur lors de l'exécution de git fetch dans le dépôt %s: %v\n", repo.Name, err)
			}

			cmd = exec.Command("git", "-C", repoDir, "for-each-ref", "--sort=-committerdate", "--count=1", "--format=%(refname:short)", "refs/heads")
			branch, err := cmd.Output()
			if err != nil {
				fmt.Printf("Erreur lors de l'obtention de la dernière branche modifiée dans le dépôt %s: %v\n", repo.Name, err)
				continue
				}

			cmd = exec.Command("git", "-C", repoDir, "pull", "origin", strings.TrimSpace(string(branch)))
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Erreur lors de l'exécution de git pull sur la dernière branche modifiée dans le dépôt %s: %v\n", repo.Name, err)
			}
		}
	// ZIP des dépôts

	// zipFile, err := os.Create("repositories.zip")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer zipFile.Close()

	// zipWriter := zip.NewWriter(zipFile)
	// defer zipWriter.Close()

	// // Ajouter des dossiers et des fichiers à l'archive ZIP
	// filepath.Walk("clones", func(path string, info os.FileInfo, err error) error {
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return err
	// 	}
	// 	header, err := zip.FileInfoHeader(info)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return err
	// 	}
	// 	header.Name = path
	// 	writer, err := zipWriter.CreateHeader(header)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return err
	// 	}
	// 	if info.IsDir() {
	// 		return nil
	// 	}
	// 	fileContent, err := os.ReadFile(path)  
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return err
	// 	}
	// 	_, err = writer.Write(fileContent)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return err
	// 	}
	// 	return nil
	// })

	http.HandleFunc("/archive", handleArchive)
	http.ListenAndServe(":8080", nil)

}

