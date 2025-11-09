package game

// Initialiser le nom du joueur
type Joueur struct {
	Name string
}

// Struct pour l'état du jeu et du board en facile
type GameEasy struct {
	Board    [6][7]int `json:"board"`
	Turn     int       `json:"turn"`
	Winner   int       `json:"winner"`
	GameOver bool      `json:"game_over"`
	Message  string    `json:"message"`
	IsDraw   bool      `json:"is_draw"`
}

// Initialise une nouvelle partie
func NewGame() *GameEasy {
	return &GameEasy{
		Turn:    1,
		Winner:  0,
		Message: "Joueur 1 (=4), à vous de jouer",
	}
}

// Ajoute un pion dans la colonne donnée
func (g *GameEasy) DropPiece(column int) bool {
	if g.GameOver {
		return false
	}

	if column < 0 || column >= 7 {
		return false // Si la colonne est invalide
	}

	// Parcourt la colonne depuis le bas
	for row := 5; row >= 0; row-- {
		if g.Board[row][column] == 0 {
			g.Board[row][column] = g.Turn

			// Vérifier la victoire
			if g.checkWinner(row, column) {
				g.GameOver = true
				g.Winner = g.Turn
				if g.Turn == 1 {
					g.Message = "<‰ Joueur 1 (=4) a gagné !"
				} else {
					g.Message = "<‰ Joueur 2 (=á) a gagné !"
				}
				return true
			}

			// Vérifier l'égalité
			if g.checkDraw() {
				g.GameOver = true
				g.IsDraw = true
				g.Message = "> Match nul ! La grille est pleine."
				return true
			}

			// Changer de joueur
			g.switchTurn()
			return true
		}
	}

	return false // Colonne pleine
}

// SwitchTurn fait passer au prochain tour de jeu
func (g *GameEasy) switchTurn() {
	if g.Turn == 1 {
		g.Turn = 2
		g.Message = "Joueur 2 (=á), à vous de jouer"
	} else {
		g.Turn = 1
		g.Message = "Joueur 1 (=4), à vous de jouer"
	}
}

// Check si un puissance 4 a été fait
func (g *GameEasy) checkWinner(row, col int) bool {
	player := g.Board[row][col]
	directions := [][2]int{
		{0, 1},  // horizontal
		{1, 0},  // vertical
		{1, 1},  // diagonale \
		{1, -1}, // diagonale /
	}

	for _, d := range directions {
		count := 1
		count += g.countDirection(row, col, d[0], d[1], player)
		count += g.countDirection(row, col, -d[0], -d[1], player)
		if count >= 4 {
			return true
		}
	}
	return false
}

// Sert pour le checkwinner, compte le nombre de même pièce à la suite
func (g *GameEasy) countDirection(row, col, dr, dc, player int) int {
	count := 0
	for {
		row += dr
		col += dc
		if row < 0 || row >= 6 || col < 0 || col >= 7 {
			break
		}
		if g.Board[row][col] == player {
			count++
		} else {
			break
		}
	}
	return count
}

// Vérifie si la grille est pleine (égalité)
func (g *GameEasy) checkDraw() bool {
	for col := 0; col < 7; col++ {
		if g.Board[0][col] == 0 {
			return false
		}
	}
	return true
}

// Reset la partie
func (g *GameEasy) Reset() {
	g.Board = [6][7]int{}
	g.Turn = 1
	g.Winner = 0
	g.GameOver = false
	g.IsDraw = false
	g.Message = "Joueur 1 (=4), à vous de jouer"
}

// GetCellClass retourne la classe CSS pour une cellule
func (g *GameEasy) GetCellClass(row, col int) string {
	switch g.Board[row][col] {
	case 1:
		return "player1"
	case 2:
		return "player2"
	default:
		return "empty"
	}
}
