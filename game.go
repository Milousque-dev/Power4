package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

//#region STRUCTURES DE DONNÉES

// Game représente une partie complète de Puissance 4
// Cette structure contient toutes les informations nécessaires pour gérer une partie
type Game struct {
	Rows           int         `json:"rows"`           // Nombre de lignes du plateau (6, 7, etc.)
	Cols           int         `json:"cols"`           // Nombre de colonnes du plateau (7, 8, 9, etc.)
	Board          [][]string  `json:"board"`          // Plateau de jeu 2D : "" = vide, "player1" ou "player2"
	CurrentPlayer  string      `json:"currentPlayer"`  // Joueur actuel ("player1" ou "player2")
	Player1        string      `json:"player1"`        // Pseudo du joueur 1
	Player2        string      `json:"player2"`        // Pseudo du joueur 2
	GameOver       bool        `json:"gameOver"`       // true si la partie est terminée
	Winner         string      `json:"winner"`         // Gagnant ("player1", "player2", "draw", ou "")
	LastMove       *Move       `json:"lastMove"`       // Dernier coup joué (nil si aucun)
	TurnCount      int         `json:"turnCount"`      // Nombre de tours joués (utilisé pour la gravité inversée)
	InverseGravity bool        `json:"inverseGravity"` // true si la gravité est actuellement inversée
}

// Move représente un coup joué sur le plateau
type Move struct {
	Row int `json:"row"` // Ligne où le jeton s'est posé (0 = haut, max = bas)
	Col int `json:"col"` // Colonne où le jeton a été déposé
}

//#endregion

//#region CRÉATION D'UNE NOUVELLE PARTIE

// NewGame crée et initialise une nouvelle partie de Puissance 4
//
// Paramètres:
//   - rows: nombre de lignes du plateau
//   - cols: nombre de colonnes du plateau
//   - player1: pseudo du premier joueur
//   - player2: pseudo du second joueur
//
// Retourne:
//   - *Game: pointeur vers la nouvelle partie initialisée
//
// La fonction initialise un plateau vide, ajoute des jetons pré-remplis selon
// la difficulté, et configure le joueur 1 comme joueur de départ
func NewGame(rows, cols int, player1, player2 string) *Game {
	// Initialisation du plateau vide (tableau 2D)
	board := make([][]string, rows)
	for i := range board {
		board[i] = make([]string, cols)
		// Remplissage avec des chaînes vides (cases vides)
		for j := range board[i] {
			board[i][j] = ""
		}
	}

	// Création de la structure Game avec les valeurs par défaut
	game := &Game{
		Rows:           rows,
		Cols:           cols,
		Board:          board,
		CurrentPlayer:  "player1", // Le joueur 1 commence toujours
		Player1:        player1,
		Player2:        player2,
		GameOver:       false,
		Winner:         "",
		LastMove:       nil,
		TurnCount:      0,
		InverseGravity: false,
	}

	// Ajout de jetons pré-remplis selon la difficulté
	numPrefilledBlocks := getPrefilledBlocksCount(rows, cols)
	game.addPrefilledBlocks(numPrefilledBlocks)

	return game
}

// getPrefilledBlocksCount détermine le nombre de jetons à pré-remplir
// selon les dimensions du plateau (qui correspondent à une difficulté)
//
// Paramètres:
//   - rows: nombre de lignes
//   - cols: nombre de colonnes
//
// Retourne:
//   - int: nombre de jetons à placer aléatoirement au départ
//
// Règles:
//   - Facile (6x7) : 3 jetons
//   - Normal (6x9) : 5 jetons
//   - Difficile (7x8) : 7 jetons
func getPrefilledBlocksCount(rows, cols int) int {
	// Mode Facile: grille 6x7 = 3 blocs pré-remplis
	if rows == 6 && cols == 7 {
		return 3
	}
	// Mode Normal: grille 6x9 = 5 blocs pré-remplis
	if rows == 6 && cols == 9 {
		return 5
	}
	// Mode Difficile: grille 7x8 = 7 blocs pré-remplis
	if rows == 7 && cols == 8 {
		return 7
	}
	// Par défaut (dimensions personnalisées), pas de blocs
	return 0
}

// addPrefilledBlocks ajoute des jetons aléatoires au plateau au démarrage
//
// Cette fonction place aléatoirement le nombre spécifié de jetons sur le plateau,
// en respectant la gravité (les jetons tombent dans les colonnes)
//
// Paramètres:
//   - count: nombre de jetons à placer
func (g *Game) addPrefilledBlocks(count int) {
	// Si count = 0, ne rien faire
	if count == 0 {
		return
	}

	// Initialisation du générateur de nombres aléatoires
	rand.Seed(time.Now().UnixNano())

	// Liste des joueurs possibles pour les jetons
	players := []string{"player1", "player2"}
	placed := 0 // Compteur de jetons placés

	// Boucle jusqu'à avoir placé tous les jetons
	for placed < count {
		// Sélection d'une colonne aléatoire
		col := rand.Intn(g.Cols)

		// Recherche de la première case vide dans la colonne (du bas vers le haut)
		row := -1
		for r := g.Rows - 1; r >= 0; r-- {
			if g.Board[r][col] == "" {
				row = r
				break
			}
		}

		// Si la colonne n'est pas pleine, placer un jeton
		if row != -1 {
			// Sélection aléatoire du joueur (player1 ou player2)
			g.Board[row][col] = players[rand.Intn(2)]
			placed++
		}
		// Si la colonne est pleine, on essaye une autre colonne au prochain tour de boucle
	}
}

//#endregion

//#region PLACEMENT DE JETONS

// DropPiece place un jeton dans une colonne spécifiée
//
// Cette fonction est le cœur de la logique de jeu :
// 1. Vérifie que le coup est valide (partie non terminée, colonne valide)
// 2. Trouve la première case vide dans la colonne (selon la gravité)
// 3. Place le jeton du joueur actuel
// 4. Incrémente le compteur de tours et gère la gravité inversée
// 5. Vérifie les conditions de victoire ou d'égalité
// 6. Change de joueur si la partie continue
//
// Paramètres:
//   - col: numéro de la colonne (0 à Cols-1)
//
// Retourne:
//   - error: nil si le coup est valide, une erreur sinon
func (g *Game) DropPiece(col int) error {
	// Vérification 1: la partie ne doit pas être terminée
	if g.GameOver {
		return errors.New("la partie est terminée")
	}

	// Vérification 2: la colonne doit être valide
	if col < 0 || col >= g.Cols {
		return fmt.Errorf("colonne invalide: %d (doit être entre 0 et %d)", col, g.Cols-1)
	}

	// Recherche de la première case vide selon la gravité actuelle
	row := -1
	if g.InverseGravity {
		// Gravité inversée : les jetons "tombent" vers le haut
		// On cherche donc la première case vide en partant du haut
		for r := 0; r < g.Rows; r++ {
			if g.Board[r][col] == "" {
				row = r
				break
			}
		}
	} else {
		// Gravité normale : les jetons tombent vers le bas
		// On cherche donc la première case vide en partant du bas
		for r := g.Rows - 1; r >= 0; r-- {
			if g.Board[r][col] == "" {
				row = r
				break
			}
		}
	}

	// Vérification 3: la colonne ne doit pas être pleine
	if row == -1 {
		return errors.New("colonne pleine")
	}

	// Placement du jeton
	g.Board[row][col] = g.CurrentPlayer
	g.LastMove = &Move{Row: row, Col: col}

	// Incrémentation du compteur de tours
	g.TurnCount++

	// Gestion de la gravité inversée : s'active tous les 5 tours
	if g.TurnCount%5 == 0 {
		g.InverseGravity = !g.InverseGravity // Inverse l'état actuel
	}

	// Vérification de la victoire (4 jetons alignés)
	if g.checkWin(row, col) {
		g.GameOver = true
		g.Winner = g.CurrentPlayer
		return nil // Partie terminée, pas besoin de changer de joueur
	}

	// Vérification de l'égalité (plateau plein)
	if g.checkDraw() {
		g.GameOver = true
		g.Winner = "draw"
		return nil // Partie terminée
	}

	// Changement de joueur pour le prochain tour
	if g.CurrentPlayer == "player1" {
		g.CurrentPlayer = "player2"
	} else {
		g.CurrentPlayer = "player1"
	}

	return nil
}

//#endregion

//#region VÉRIFICATION DE VICTOIRE

// checkWin vérifie si le dernier coup joué a créé un alignement gagnant
//
// Un joueur gagne en alignant 4 jetons dans l'une des 4 directions :
// - Horizontal (←→)
// - Vertical (↑↓)
// - Diagonale descendante (\)
// - Diagonale montante (/)
//
// Paramètres:
//   - row: ligne du dernier jeton placé
//   - col: colonne du dernier jeton placé
//
// Retourne:
//   - bool: true si le coup est gagnant, false sinon
func (g *Game) checkWin(row, col int) bool {
	// Récupération du joueur qui vient de jouer
	player := g.Board[row][col]

	// Définition des 4 directions à vérifier
	// Chaque direction est définie par un vecteur [dRow, dCol]
	directions := [][2]int{
		{0, 1},  // Horizontal (même ligne, colonnes ++)
		{1, 0},  // Vertical (lignes ++, même colonne)
		{1, 1},  // Diagonale descendante (\)
		{1, -1}, // Diagonale montante (/)
	}

	// Vérification dans chaque direction
	for _, dir := range directions {
		// Comptage des jetons alignés :
		// 1 (le jeton actuel) +
		// count dans une direction +
		// count dans la direction opposée
		count := 1 +
			g.countDirection(row, col, dir[0], dir[1], player) +     // Direction positive
			g.countDirection(row, col, -dir[0], -dir[1], player)     // Direction négative (opposée)

		// Si on trouve 4 jetons alignés ou plus : victoire !
		if count >= 4 {
			return true
		}
	}

	// Aucun alignement de 4 jetons trouvé
	return false
}

// countDirection compte le nombre de jetons alignés dans une direction donnée
//
// Cette fonction part de la position actuelle et compte combien de jetons
// du même joueur sont alignés dans la direction spécifiée, jusqu'à rencontrer
// une case vide, un jeton adverse, ou le bord du plateau
//
// Paramètres:
//   - row: ligne de départ
//   - col: colonne de départ
//   - dRow: déplacement en ligne (+1 = bas, -1 = haut, 0 = même ligne)
//   - dCol: déplacement en colonne (+1 = droite, -1 = gauche, 0 = même colonne)
//   - player: joueur dont on compte les jetons
//
// Retourne:
//   - int: nombre de jetons alignés dans cette direction (sans compter la position de départ)
func (g *Game) countDirection(row, col, dRow, dCol int, player string) int {
	count := 0
	// Position suivante dans la direction
	r, c := row+dRow, col+dCol

	// Continue tant qu'on est dans les limites ET que la case contient le même joueur
	for r >= 0 && r < g.Rows && c >= 0 && c < g.Cols && g.Board[r][c] == player {
		count++
		// Déplacement vers la case suivante
		r += dRow
		c += dCol
	}

	return count
}

//#endregion

//#region VÉRIFICATION D'ÉGALITÉ

// checkDraw vérifie si le plateau est complètement rempli (égalité)
//
// Une égalité se produit quand toutes les cases sont remplies
// et qu'aucun joueur n'a aligné 4 jetons
//
// Retourne:
//   - bool: true si le plateau est plein, false s'il reste au moins une case vide
func (g *Game) checkDraw() bool {
	// Parcours de toutes les cases du plateau
	for row := 0; row < g.Rows; row++ {
		for col := 0; col < g.Cols; col++ {
			// Si on trouve une case vide, ce n'est pas une égalité
			if g.Board[row][col] == "" {
				return false
			}
		}
	}
	// Toutes les cases sont remplies : c'est une égalité
	return true
}

//#endregion

//#region EXPORT DE L'ÉTAT

// GetState retourne l'état complet du jeu sous forme de map
//
// Cette fonction est utilisée pour sérialiser l'état du jeu en JSON
// et l'envoyer au client frontend
//
// Retourne:
//   - map[string]interface{}: dictionnaire contenant toutes les informations de la partie
func (g *Game) GetState() map[string]interface{} {
	return map[string]interface{}{
		"rows":           g.Rows,
		"cols":           g.Cols,
		"board":          g.Board,
		"currentPlayer":  g.CurrentPlayer,
		"player1":        g.Player1,
		"player2":        g.Player2,
		"gameOver":       g.GameOver,
		"winner":         g.Winner,
		"lastMove":       g.LastMove,
		"turnCount":      g.TurnCount,
		"inverseGravity": g.InverseGravity,
	}
}

//#endregion
