# Golang_CC

Énoncé 4A IBC : Langage GO
Projet : Récupération massive de Github Repositories
Descriptif :
Ce projet a pour but de tester les compétences des étudiants de 4ème année à créer une application
en utilisant leurs connaissances sur le langage Go.
Modalités :
Ce projet est à rendre individuellement.
Le rendu doit être sous le format d’un repository Git public.
Le rendu doit être uploadé sur Github avant la date butoir du projet.
Enoncé:
Le but de ce projet annuel est de créer une application en Golang complète permettant de cloner les
repositories depuis Github, selon les critères ci-dessous.
Les étudiants seront amenés à développer les fonctionnalités suivantes :
- Créer une application qui requête l’API Github pour récupérer:
- la liste de repositories d’un utilisateur,
- ou la liste de repositories d’une organisation,
- Trier ces repositories par dernière modification.
- L’application doit récupérer au minimum TOUS les repositories spécifiés, ou au minimum les
100 derniers repositories par date de modification.
- L’application doit écrire un CSV de cette liste, avec l’ensemble des informations récupérées
sur l’API.
- L’application doit cloner ces repositories en local.
- L’application doit exécuter un Git Pull sur la dernière branche modifiée (dernier commit) en
local.
- L’application doit aussi exécuter un Git Fetch pour récupérer toutes les références de
branches en local.
- L’application doit créer une archive (ZIP ou 7z) de ces repositories à la fin du traitement en
local.
Une fois déployée, la dApp aura comme fonctionnalités de :
- Spécifier le pseudo Github d’un utilisateur ou une organisation,
- Lister et cloner les repositories publiques de l’utilisateur ou l’organisation,
- Si un Token API Github est fourni, l’application doit en supplément cloner les repositories
privés de l’utilisateur ou l’organisation,
- Rendre disponible le téléchargement de ces repositories via une API.
L’utilisation des notions Golang suivantes est obligatoire :
- Webserver HTTP pour le téléchargement de l’archive,
- Goroutines & Waitgroups pour optimiser l’exécution du code.
La dApp doit être Dockerisée afin de faciliter son déploiement.
Des volumes persistants pour la BDD sont à prévoir.
Annexes :
https://pkg.go.dev/
https://docs.github.com/fr/rest?apiVersion=2022-11-28


    # GitHub Repository Cloner and Archiver

Ce projet permet de cloner des dépôts GitHub d'un utilisateur ou d'une organisation spécifique, de mettre à jour les dépôts avec la dernière branche modifiée, et finalement de créer une archive ZIP contenant tous les dépôts clonés.

## Prérequis

- Go 1.18 ou version ultérieure (en raison de l'utilisation de `io` et `os` pour la lecture des fichiers)
- Un compte GitHub et un [token d'accès personnel](https://github.com/settings/tokens) pour l'authentification

## Configuration

1. Définissez les variables d'environnement suivantes:
   - `GITHUB_ORG` : Le nom de l'organisation GitHub (si vous voulez cloner les dépôts d'une organisation)
   - `GITHUB_USER` : Le nom d'utilisateur GitHub (si vous voulez cloner les dépôts d'un utilisateur)
   - `GITHUB_TOKEN` : Votre token d'accès personnel GitHub

> **Note**: Vous pouvez définir soit `GITHUB_ORG` soit `GITHUB_USER`, pas les deux. 

## Comment ça marche

1. Le programme fait d'abord une requête à l'API GitHub pour obtenir une liste des dépôts de l'utilisateur ou de l'organisation spécifiée.
2. Il trie les dépôts par date de dernière modification.
3. Il crée un fichier CSV contenant des informations sur chaque dépôt.
4. Il clone chaque dépôt dans un dossier nommé `clones`.
5. Il exécute `git fetch` et `git pull` sur la dernière branche modifiée de chaque dépôt cloné.
6. Enfin, il crée une archive ZIP nommée `repositories.zip` contenant tous les dépôts clonés.

## Usage

1. Assurez-vous d'avoir défini les variables d'environnement nécessaires.
2. Exécutez le programme en utilisant la commande `go run main.go` dans le répertoire du projet.
3. Une fois le programme terminé, vous trouverez l'archive ZIP `repositories.zip` dans le répertoire du projet, et un fichier CSV `repositories.csv` contenant des informations sur chaque dépôt.

## Notes Importantes

- L'archive ZIP contiendra un dossier `clones` qui contient tous les dépôts clonés.
- Le fichier CSV contiendra les colonnes suivantes : `Name`, `Description`, `Updated At`, et `Clone URL` pour chaque dépôt.
- Assurez-vous que votre token d'accès personnel GitHub a les autorisations nécessaires pour cloner les dépôts.

Ce projet sert d'exemple basique et pourrait être étendu ou modifié pour répondre à des besoins spécifiques.
