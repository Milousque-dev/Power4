/**
 * PUISSANCE 4 - SERVEUR HTTP PRINCIPAL
 *
 * Ce fichier contient le point d'entrée de l'application serveur.
 * Il configure le serveur HTTP Go, les routes, et démarre l'écoute
 * sur le port 8080.
 *
 * Architecture:
 *   - Serveur de fichiers statiques (CSS, JS, images)
 *   - 3 pages HTML (difficulté, skins, jeu)
 *   - API REST pour la logique du jeu
 */

package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

//#region POINT D'ENTRÉE PRINCIPAL

// main est la fonction principale du serveur
//
// Cette fonction:
// 1. Crée un gestionnaire de parties
// 2. Configure les routes pour les fichiers statiques
// 3. Configure les routes pour les pages HTML
// 4. Configure les routes de l'API
// 5. Démarre le serveur HTTP
func main() {
	//#region Initialisation

	// Création du gestionnaire de parties
	// Il maintiendra l'état de la partie en cours
	gameManager := NewGameManager()

	//#endregion

	//#region Configuration des routes - Fichiers statiques

	// Serveur de fichiers CSS (styles)
	// Route: /css/* → dossier /css
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	// Serveur de fichiers JavaScript (logique client)
	// Route: /js/* → dossier /js
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	// Serveur de fichiers statiques (images de jetons, fonds d'écran)
	// Route: /static/* → dossier /static
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//#endregion

	//#region Configuration des routes - Pages HTML

	// Page d'accueil: Sélection de la difficulté
	// Route: GET /
	// Template: templates/difficulty.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// S'assurer que seule la racine est traitée ici
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// Chargement et exécution du template
		tmpl, err := template.ParseFiles("templates/difficulty.html")
		if err != nil {
			http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
			log.Println("Erreur template:", err)
			return
		}
		tmpl.Execute(w, nil)
	})

	// Page de sélection des skins et pseudos
	// Route: GET /skins
	// Template: templates/skins.html
	http.HandleFunc("/skins", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/skins.html")
		if err != nil {
			http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
			log.Println("Erreur template:", err)
			return
		}
		tmpl.Execute(w, nil)
	})

	// Page de jeu
	// Route: GET /game
	// Template: templates/game.html
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/game.html")
		if err != nil {
			http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
			log.Println("Erreur template:", err)
			return
		}
		tmpl.Execute(w, nil)
	})

	//#endregion

	//#region Configuration des routes - API REST

	// API: Créer une nouvelle partie
	// Route: POST /api/game/new
	// Body: {rows, cols, player1, player2}
	// Réponse: État initial de la partie
	http.HandleFunc("/api/game/new", gameManager.HandleNewGame)

	// API: Placer un jeton
	// Route: POST /api/game/drop
	// Body: {col}
	// Réponse: Nouvel état de la partie
	http.HandleFunc("/api/game/drop", gameManager.HandleDropPiece)

	// API: Obtenir l'état actuel
	// Route: GET /api/game/state
	// Réponse: État actuel de la partie
	http.HandleFunc("/api/game/state", gameManager.HandleGetState)

	// API: Réinitialiser le jeu
	// Route: POST /api/game/reset
	// Réponse: Message de confirmation
	http.HandleFunc("/api/game/reset", gameManager.HandleReset)

	//#endregion

	//#region Démarrage du serveur

	// Récupération du port depuis la variable d'environnement PORT
	// Par défaut: 8080 si la variable n'est pas définie
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Log du démarrage
	log.Printf("🎮 Serveur Puissance 4 démarré sur http://localhost:%s\n", port)

	// Démarrage du serveur HTTP
	// Bloque jusqu'à une erreur fatale ou interruption
	log.Fatal(http.ListenAndServe(":"+port, nil))

	//#endregion
}

//#endregion

//#region FONCTIONS UTILITAIRES HTTP

// respondJSON envoie une réponse JSON au client
//
// Cette fonction utilitaire simplifie l'envoi de réponses JSON
// en gérant automatiquement:
// - Le header Content-Type
// - Le code de statut HTTP
// - La sérialisation en JSON
//
// Paramètres:
//   - w: ResponseWriter pour envoyer la réponse
//   - status: Code de statut HTTP (200, 400, 404, etc.)
//   - data: Données à sérialiser en JSON (struct, map, etc.)
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	// Définition du type de contenu
	w.Header().Set("Content-Type", "application/json")

	// Définition du code de statut
	w.WriteHeader(status)

	// Sérialisation et envoi des données
	json.NewEncoder(w).Encode(data)
}

// respondError envoie une erreur JSON au client
//
// Cette fonction est un raccourci pour envoyer une erreur
// formatée de manière cohérente.
//
// Paramètres:
//   - w: ResponseWriter pour envoyer la réponse
//   - status: Code d'erreur HTTP (400, 404, 500, etc.)
//   - message: Message d'erreur descriptif
//
// Format de réponse JSON:
//
//	{
//	  "error": "message d'erreur"
//	}
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

//#endregion
