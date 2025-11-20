package main

import (
	"encoding/json"
	"net/http"
)

//#region GESTIONNAIRE DE PARTIES

// GameManager gère l'état d'une partie en cours
//
// Cette structure maintient une seule partie active à la fois.
// Elle pourrait être étendue pour gérer plusieurs parties simultanées
// (par exemple avec un map[sessionID]*Game)
type GameManager struct {
	game *Game // Partie actuellement en cours (nil si aucune partie active)
}

// NewGameManager crée un nouveau gestionnaire de jeu
//
// Retourne:
//   - *GameManager: nouveau gestionnaire sans partie active
func NewGameManager() *GameManager {
	return &GameManager{
		game: nil, // Aucune partie au démarrage
	}
}

//#endregion

//#region HANDLERS HTTP - CRÉATION DE PARTIE

// HandleNewGame gère la création d'une nouvelle partie
//
// Route: POST /api/game/new
// Body JSON attendu:
//
//	{
//	  "rows": 6,
//	  "cols": 7,
//	  "player1": "Alice",
//	  "player2": "Bob"
//	}
//
// Paramètres:
//   - w: ResponseWriter pour envoyer la réponse
//   - r: Request contenant les données de la nouvelle partie
//
// Réponse:
//   - 200 OK: État initial de la partie
//   - 400 Bad Request: Paramètres invalides
//   - 405 Method Not Allowed: Méthode HTTP incorrecte
func (gm *GameManager) HandleNewGame(w http.ResponseWriter, r *http.Request) {
	// Vérification de la méthode HTTP
	if r.Method != "POST" {
		respondError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	// Structure pour décoder le JSON de la requête
	var req struct {
		Rows    int    `json:"rows"`    // Nombre de lignes souhaité
		Cols    int    `json:"cols"`    // Nombre de colonnes souhaité
		Player1 string `json:"player1"` // Pseudo du joueur 1
		Player2 string `json:"player2"` // Pseudo du joueur 2
	}

	// Décodage du JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}

	// Validation: nombre de lignes (4 minimum pour pouvoir gagner, 10 maximum pour l'UI)
	if req.Rows < 4 || req.Rows > 10 {
		respondError(w, http.StatusBadRequest, "Nombre de lignes doit être entre 4 et 10")
		return
	}

	// Validation: nombre de colonnes (4 minimum pour pouvoir gagner, 10 maximum pour l'UI)
	if req.Cols < 4 || req.Cols > 10 {
		respondError(w, http.StatusBadRequest, "Nombre de colonnes doit être entre 4 et 10")
		return
	}

	// Validation: les pseudos ne doivent pas être vides
	if req.Player1 == "" || req.Player2 == "" {
		respondError(w, http.StatusBadRequest, "Les pseudos sont obligatoires")
		return
	}

	// Création de la nouvelle partie
	gm.game = NewGame(req.Rows, req.Cols, req.Player1, req.Player2)

	// Envoi de l'état initial au client
	respondJSON(w, http.StatusOK, gm.game.GetState())
}

//#endregion

//#region HANDLERS HTTP - PLACEMENT DE JETON

// HandleDropPiece gère le placement d'un jeton dans une colonne
//
// Route: POST /api/game/drop
// Body JSON attendu:
//
//	{
//	  "col": 3
//	}
//
// Paramètres:
//   - w: ResponseWriter pour envoyer la réponse
//   - r: Request contenant le numéro de colonne
//
// Réponse:
//   - 200 OK: Nouvel état de la partie après le coup
//   - 400 Bad Request: Coup invalide ou aucune partie en cours
//   - 405 Method Not Allowed: Méthode HTTP incorrecte
func (gm *GameManager) HandleDropPiece(w http.ResponseWriter, r *http.Request) {
	// Vérification de la méthode HTTP
	if r.Method != "POST" {
		respondError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	// Vérification qu'une partie est en cours
	if gm.game == nil {
		respondError(w, http.StatusBadRequest, "Aucune partie en cours")
		return
	}

	// Structure pour décoder le JSON de la requête
	var req struct {
		Col int `json:"col"` // Numéro de colonne où placer le jeton
	}

	// Décodage du JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Format JSON invalide")
		return
	}

	// Tentative de placement du jeton
	err := gm.game.DropPiece(req.Col)
	if err != nil {
		// Erreur de jeu (colonne pleine, partie terminée, etc.)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Envoi du nouvel état au client
	respondJSON(w, http.StatusOK, gm.game.GetState())
}

//#endregion

//#region HANDLERS HTTP - CONSULTATION ET RÉINITIALISATION

// HandleGetState retourne l'état actuel du jeu
//
// Route: GET /api/game/state
//
// Paramètres:
//   - w: ResponseWriter pour envoyer la réponse
//   - r: Request
//
// Réponse:
//   - 200 OK: État actuel de la partie
//   - 400 Bad Request: Aucune partie en cours
//   - 405 Method Not Allowed: Méthode HTTP incorrecte
func (gm *GameManager) HandleGetState(w http.ResponseWriter, r *http.Request) {
	// Vérification de la méthode HTTP
	if r.Method != "GET" {
		respondError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	// Vérification qu'une partie existe
	if gm.game == nil {
		respondError(w, http.StatusBadRequest, "Aucune partie en cours")
		return
	}

	// Envoi de l'état actuel
	respondJSON(w, http.StatusOK, gm.game.GetState())
}

// HandleReset réinitialise le jeu (supprime la partie en cours)
//
// Route: POST /api/game/reset
//
// Paramètres:
//   - w: ResponseWriter pour envoyer la réponse
//   - r: Request
//
// Réponse:
//   - 200 OK: Message de confirmation
//   - 405 Method Not Allowed: Méthode HTTP incorrecte
func (gm *GameManager) HandleReset(w http.ResponseWriter, r *http.Request) {
	// Vérification de la méthode HTTP
	if r.Method != "POST" {
		respondError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	// Suppression de la partie en cours
	gm.game = nil

	// Confirmation de la réinitialisation
	respondJSON(w, http.StatusOK, map[string]string{"message": "Jeu réinitialisé"})
}

//#endregion
