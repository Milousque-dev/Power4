# 📚 Documentation Technique - Puissance 4

## Table des matières
1. [Architecture globale](#architecture-globale)
2. [Flux de données](#flux-de-données)
3. [Modules détaillés](#modules-détaillés)
4. [Algorithmes clés](#algorithmes-clés)
5. [Guide du débutant](#guide-du-débutant)

---

## Architecture globale

### Séparation des responsabilités

Le projet suit une architecture **Client-Serveur** avec séparation stricte entre backend et frontend :

```
┌──────────────────────────────────────────────────────────┐
│                    BACKEND (Go)                          │
│  ┌──────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ main.go  │→ │game_manager  │→ │   game.go    │      │
│  │  HTTP    │  │   .go        │  │   Logique    │      │
│  │ Serveur  │  │  Gestionnaire│  │   du jeu     │      │
│  └──────────┘  └──────────────┘  └──────────────┘      │
│       ↕ API REST (JSON)                                  │
└──────────────────────────────────────────────────────────┘
                        ↕
┌──────────────────────────────────────────────────────────┐
│                  FRONTEND (HTML/CSS/JS)                  │
│  ┌─────────────┐                                         │
│  │ index.html  │  Structure & Présentation               │
│  └──────┬──────┘                                         │
│         │                                                 │
│    ┌────┴────┬──────────┬──────────┐                    │
│    │         │          │          │                     │
│ ┌──▼─────┐ ┌▼──────┐ ┌─▼──────┐ ┌─▼──────────┐         │
│ │styles  │ │game   │ │ui      │ │fireworks   │         │
│ │.css    │ │.js    │ │.js     │ │.js         │         │
│ └────────┘ └───────┘ └────────┘ └────────────┘         │
│   Style    API Calls  Interface   Animations            │
└──────────────────────────────────────────────────────────┘
```

**IMPORTANT** : La logique du jeu (placement des jetons, détection de victoire, égalité) est **entièrement gérée par le backend Go**. JavaScript ne fait que communiquer avec l'API et mettre à jour l'interface.

### Ordre de chargement

**Important** : Les scripts doivent être chargés dans cet ordre :

1. `game.js` - Déclare les variables globales
2. `fireworks.js` - Utilise ces variables
3. `ui.js` - Utilise les fonctions de game.js

---

## Flux de données

### 1. Démarrage du jeu

```
Chargement page
    ↓
Affichage écran difficulté
    ↓
Clic sur difficulté
    ↓
selectDifficulty(difficulty)
    ├─ Définit ROWS et COLS
    ├─ Change le thème visuel
    └─ Affiche écran sélection skins
```

### 2. Configuration des joueurs

```
Écran sélection skins
    ↓
Saisie pseudos + sélection skins
    ↓
checkIfReady()
    ├─ Vérifie pseudos remplis
    ├─ Vérifie skins différents
    └─ Active/désactive bouton démarrage
    ↓
Clic "Commencer"
    ↓
initBoard()
    └─ Crée la grille de jeu
```

### 3. Déroulement d'une partie

```
Tour de jeu
    ↓
Clic sur colonne
    ↓
dropPiece(col)
    ├─ Trouve case vide la plus basse
    ├─ Place le jeton (board[row][col])
    ├─ Lance animation chute
    ├─ Change joueur actuel
    └─ Après 600ms:
        ├─ checkWin(row, col)
        │   └─ Si gagnant → createFireworks()
        └─ checkDraw()
            └─ Si égalité → Affiche message
```

---

## Modules détaillés

### game.js - Logique du jeu

#### Variables globales

```javascript
// Dimensions de la grille (varient selon difficulté)
let ROWS = 6;  // Nombre de lignes
let COLS = 7;  // Nombre de colonnes

// État du jeu
let currentDifficulty = '';  // 'easy', 'normal', 'hard'
let currentPlayer = 'player1';  // 'player1' ou 'player2'
let board = [];  // Tableau 2D: board[ligne][colonne]
let gameOver = false;  // true si partie terminée

// Configuration des joueurs
let selectedSkins = {
    player1: 'skinX',  // Nom du fichier du skin
    player2: 'skinY'
};
let playerPseudos = {
    player1: 'Pseudo1',
    player2: 'Pseudo2'
};
```

#### Fonction initBoard()

**Rôle** : Crée le plateau de jeu visuel et le tableau de données

```javascript
function initBoard() {
    // 1. Créer tableau 2D vide
    board = Array(ROWS).fill(null).map(() => Array(COLS).fill(null));

    // 2. Réinitialiser état
    gameOver = false;
    currentPlayer = 'player1';

    // 3. Créer grille HTML
    for (let col = 0; col < COLS; col++) {
        for (let row = 0; row < ROWS; row++) {
            // Créer cellule avec ID unique
            cell.id = `cell-${row}-${col}`;
        }
    }
}
```

**Exemple de grille 3×3** :

```
HTML:                      Tableau board:
┌─────┬─────┬─────┐       [[null, null, null],
│ 0,0 │ 0,1 │ 0,2 │        [null, null, null],
├─────┼─────┼─────┤        [null, null, null]]
│ 1,0 │ 1,1 │ 1,2 │
├─────┼─────┼─────┤
│ 2,0 │ 2,1 │ 2,2 │
└─────┴─────┴─────┘
```

#### Fonction dropPiece(col)

**Rôle** : Place un jeton dans une colonne

**Algorithme** :
1. Parcourir la colonne du bas vers le haut
2. Trouver la première case vide (null)
3. Placer le jeton à cette position
4. Lancer l'animation
5. Changer de joueur
6. Vérifier victoire/égalité après 600ms

```javascript
// Exemple: Placer un jeton en colonne 2
dropPiece(2);

// Avant:                Après:
// [null, null, null]    [null, null, null]
// [null, null, null]    [null, null, null]
// [null, null, null]    [null, null, 'player1']
//                              ↑
//                        Jeton placé ici
```

#### Algorithme de détection de victoire

**checkWin(row, col)** vérifie 4 directions depuis le dernier jeton placé :

```
        ↖ diagonal \    ↑ vertical
             ↖          ↑
              ↖         ↑
               ●────────→  horizontal
              ↗         ↓
             ↗          ↓
        ↗ diagonal /    ↓
```

**Principe** :
- On compte dans une direction et dans la direction opposée
- Si total ≥ 3 (+ le jeton actuel = 4), c'est gagné !

```javascript
// Exemple horizontal:
// Jeton placé à (2, 3)
//
// [X, X, X, ●, X, X, X]
//           ↑
//        Position actuelle
//
// Compte à gauche: 3 jetons (0,1,2)
// Compte à droite: 0 jetons
// Total: 3 (+ actuel) = 4 → VICTOIRE!
```

---

### ui.js - Interface utilisateur

#### Fonction selectDifficulty(difficulty)

**Configuration des grilles** :

| Difficulté | Lignes | Colonnes | Total cases |
|------------|--------|----------|-------------|
| easy       | 6      | 7        | 42          |
| normal     | 7      | 8        | 56          |
| hard       | 8      | 9        | 72          |

#### Gestion des skins

**Problème** : Les deux joueurs ne doivent pas avoir le même jeton

**Solution** : `updateSkinAvailability()`

```javascript
// État initial: tous les skins disponibles
[skin1] [skin2] [skin3] [skin4]

// Joueur 1 choisit skin2
[skin1] [skin2*] [skin3] [skin4]

// Pour Joueur 2, skin2 devient disabled
[skin1] [skin2-gris] [skin3] [skin4]
         ↑ Ne peut plus être sélectionné
```

**Classes CSS appliquées** :
- `.selected` : Bordure bleue, sélectionné par le joueur
- `.disabled` : Grisé, non cliquable (déjà pris)

---

### fireworks.js - Animation

#### Structure d'un feu d'artifice

```
     Explosion
        ●
    ╱   │   ╲
   ╱    │    ╲
  ╱     │     ╲
 ●      ●      ●
╱       │       ╲
        │
     50 particules
   disposées en cercle
```

#### Calcul de la position des particules

**Principe** : Répartition circulaire avec trigonométrie

```javascript
// Pour chaque particule i (0 à 49):
angle = (2π × i) / 50

// Position X et Y basées sur le cercle trigonométrique
x = cos(angle) × vitesse
y = sin(angle) × vitesse
```

**Exemple visuel** :

```
     i=12        i=0
       ●   ...   ●
   i=25 ●       ● i=37

   i=37 ●   🎆   ● i=12
       ●   ...   ●
     i=49        i=25
```

#### Timeline de l'animation

```
0ms    ├─ Créer 15 feux (250ms d'intervalle)
       │
3750ms ├─ Dernier feu initial
       │
4000ms ├─ Intervalle continu (400ms)
       │  ├─ Nouveau feu
       │  ├─ Nouveau feu
       │  └─ Nouveau feu
       │
5000ms └─ Arrêt de l'animation
```

---

## Algorithmes clés

### 1. Détection d'alignement (checkDirection)

**Complexité** : O(n) où n est la taille de la grille

```javascript
function checkDirection(row, col, dRow, dCol, player) {
    let count = 0;
    let r = row + dRow;
    let c = col + dCol;

    while (estDansGrille(r, c) && board[r][c] === player) {
        count++;
        r += dRow;
        c += dCol;
    }

    return count;
}
```

**Exemple horizontal (dRow=0, dCol=1)** :

```
Position départ: (2, 3)
Direction: (0, 1) → Vers la droite

Itération 1: (2, 4) → board[2][4] === player? Oui → count=1
Itération 2: (2, 5) → board[2][5] === player? Oui → count=2
Itération 3: (2, 6) → board[2][6] === player? Non → Stop
Résultat: count=2
```

### 2. Génération de grille (initBoard)

**Complexité** : O(rows × cols)

```javascript
// Pour une grille 6×7:
// 6 colonnes × 7 lignes = 42 itérations
for (let col = 0; col < COLS; col++) {      // 7 fois
    for (let row = 0; row < ROWS; row++) {  // 6 fois
        // Créer cellule
    }
}
// Total: 7 × 6 = 42 créations d'éléments DOM
```

### 3. Vérification d'égalité (checkDraw)

**Complexité** : O(rows × cols) dans le pire cas

```javascript
function checkDraw() {
    for (let row = 0; row < ROWS; row++) {
        for (let col = 0; col < COLS; col++) {
            if (board[row][col] === null) {
                return false;  // Case vide trouvée
            }
        }
    }
    return true;  // Aucune case vide
}
```

**Optimisation** : S'arrête dès qu'une case vide est trouvée

---

## Guide du débutant

### Comment lire le code JavaScript

#### 1. Variables

```javascript
// Déclaration avec let (valeur modifiable)
let score = 0;

// Déclaration avec const (valeur fixe)
const MAX_PLAYERS = 2;

// Tableau (liste d'éléments)
let jetons = ['rouge', 'jaune', 'bleu'];

// Objet (collection de propriétés)
let joueur = {
    nom: 'Pierre',
    score: 10
};
```

#### 2. Fonctions

```javascript
// Déclaration
function direBonjour(nom) {
    return "Bonjour " + nom;
}

// Appel
let message = direBonjour("Marie");  // "Bonjour Marie"
```

#### 3. Conditions

```javascript
if (score > 100) {
    console.log("Bravo !");
} else if (score > 50) {
    console.log("Bien !");
} else {
    console.log("Continue !");
}
```

#### 4. Boucles

```javascript
// Boucle for (nombre fixe d'itérations)
for (let i = 0; i < 5; i++) {
    console.log(i);  // Affiche: 0, 1, 2, 3, 4
}

// Boucle while (tant que condition vraie)
while (score < 100) {
    score += 10;
}
```

#### 5. DOM (Document Object Model)

```javascript
// Récupérer un élément par son ID
let bouton = document.getElementById('monBouton');

// Modifier le contenu texte
bouton.textContent = "Cliquez ici";

// Ajouter une classe CSS
bouton.classList.add('actif');

// Créer un nouvel élément
let div = document.createElement('div');
```

### Exercices pratiques

#### Exercice 1: Ajouter un compteur de coups

```javascript
// Dans game.js, ajouter:
let coupJoues = 0;

// Dans dropPiece(), après avoir placé un jeton:
coupJoues++;
console.log("Coups joués:", coupJoues);
```

#### Exercice 2: Changer la couleur du message de victoire

```css
/* Dans styles.css: */
.message.winner {
    background: #ff6347;  /* Rouge tomate */
    color: white;
}
```

#### Exercice 3: Ajouter un son de victoire

```javascript
// Dans dropPiece(), après "Bien vue...":
let son = new Audio('chemin/vers/victoire.mp3');
son.play();
```

---

## Débogage

### Outils de développement (F12)

#### Console
```javascript
// Afficher des valeurs pour déboguer
console.log("Valeur de ROWS:", ROWS);
console.log("Plateau:", board);
console.log("Joueur actuel:", currentPlayer);
```

#### Breakpoints
1. Ouvrir l'onglet "Sources"
2. Cliquer sur le numéro de ligne
3. Le code s'arrêtera à cet endroit
4. Inspecter les variables

### Erreurs courantes

#### 1. "Cannot read property of null"
```javascript
// Problème:
let element = document.getElementById('mauvaisID');
element.textContent = "Test";  // ❌ element est null

// Solution:
let element = document.getElementById('bonID');
if (element) {  // ✅ Vérifier d'abord
    element.textContent = "Test";
}
```

#### 2. Variables non définies
```javascript
// Problème:
console.log(maVariable);  // ❌ ReferenceError

// Solution:
let maVariable = 5;
console.log(maVariable);  // ✅
```

#### 3. Boucle infinie
```javascript
// Problème:
let i = 0;
while (i < 10) {
    console.log(i);
    // ❌ Oubli d'incrémenter i
}

// Solution:
let i = 0;
while (i < 10) {
    console.log(i);
    i++;  // ✅ Incrémenter
}
```

---

## Backend Go - Architecture détaillée

### Fichiers du backend

#### 1. `main.go` - Serveur HTTP
```go
func main() {
    gameManager := NewGameManager()

    // Serveur de fichiers statiques
    fs := http.FileServer(http.Dir("."))
    http.Handle("/", fs)

    // Routes API
    http.HandleFunc("/api/game/new", corsMiddleware(gameManager.HandleNewGame))
    http.HandleFunc("/api/game/drop", corsMiddleware(gameManager.HandleDropPiece))
    http.HandleFunc("/api/game/state", corsMiddleware(gameManager.HandleGetState))
    http.HandleFunc("/api/game/reset", corsMiddleware(gameManager.HandleReset))

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**Rôle** : Point d'entrée du serveur, configure les routes API et sert les fichiers statiques.

#### 2. `game.go` - Logique du jeu

**Structure Game** :
```go
type Game struct {
    Rows          int         // Nombre de lignes
    Cols          int         // Nombre de colonnes
    Board         [][]string  // Plateau (tableau 2D)
    CurrentPlayer string      // "player1" ou "player2"
    GameOver      bool        // Partie terminée?
    Winner        string      // "", "player1", "player2", "draw"
    LastMove      *Move       // Dernier coup joué
}
```

**Fonctions principales** :
- `NewGame()` : Crée une nouvelle partie avec un plateau vide
- `DropPiece(col)` : Place un jeton dans une colonne
  1. Vérifie que la colonne n'est pas pleine
  2. Trouve la première case vide en partant du bas
  3. Place le jeton
  4. Vérifie s'il y a victoire ou égalité
  5. Change de joueur si la partie continue

- `checkWin(row, col)` : Vérifie si le dernier coup est gagnant
  - Vérifie 4 directions : horizontal, vertical, 2 diagonales
  - Pour chaque direction, compte dans les deux sens
  - Retourne `true` si 4 jetons alignés ou plus

- `checkDraw()` : Vérifie si le plateau est plein (égalité)

**Algorithme de détection de victoire** :
```go
directions := [][2]int{
    {0, 1},   // Horizontal →
    {1, 0},   // Vertical ↓
    {1, 1},   // Diagonale ↘
    {1, -1},  // Diagonale ↙
}

for _, dir := range directions {
    count := 1 + // Jeton actuel
        countDirection(row, col, dir[0], dir[1], player) +      // Direction positive
        countDirection(row, col, -dir[0], -dir[1], player)      // Direction négative

    if count >= 4 {
        return true // Victoire!
    }
}
```

#### 3. `game_manager.go` - Gestionnaire de parties

**Structure GameManager** :
```go
type GameManager struct {
    mu   sync.RWMutex // Mutex pour accès concurrent
    game *Game        // Partie en cours
}
```

**Handlers HTTP** :
- `HandleNewGame` : POST /api/game/new
  - Reçoit : `{rows, cols, player1, player2}`
  - Crée une nouvelle partie
  - Retourne l'état initial

- `HandleDropPiece` : POST /api/game/drop
  - Reçoit : `{col}`
  - Appelle `game.DropPiece(col)`
  - Retourne le nouvel état

- `HandleGetState` : GET /api/game/state
  - Retourne l'état actuel du jeu

- `HandleReset` : POST /api/game/reset
  - Réinitialise la partie

**Sécurité concurrente** :
Le `sync.RWMutex` permet à plusieurs lecteurs d'accéder en même temps à l'état du jeu, mais un seul écrivain peut modifier l'état à la fois.

### API REST - Endpoints

| Méthode | Endpoint | Body | Réponse |
|---------|----------|------|---------|
| POST | `/api/game/new` | `{rows, cols, player1, player2}` | État initial du jeu |
| POST | `/api/game/drop` | `{col}` | Nouvel état après le coup |
| GET | `/api/game/state` | - | État actuel |
| POST | `/api/game/reset` | - | `{message: "Jeu réinitialisé"}` |

### Communication Frontend-Backend

**Flux d'un coup** :
```
1. Joueur clique sur colonne
   ↓
2. JavaScript appelle callAPI('/game/drop', 'POST', {col: 3})
   ↓
3. Requête HTTP vers http://localhost:8080/api/game/drop
   ↓
4. Backend Go traite la demande:
   - Valide le coup
   - Met à jour l'état du jeu
   - Vérifie victoire/égalité
   ↓
5. Backend répond avec JSON:
   {
     board: [...],
     currentPlayer: "player2",
     gameOver: false,
     winner: "",
     lastMove: {row: 5, col: 3}
   }
   ↓
6. JavaScript met à jour l'UI:
   - Affiche le jeton à la position retournée
   - Lance l'animation de chute
   - Met à jour l'indicateur de joueur
   - Affiche le message de victoire si gameOver
```

### Lancer le serveur

```bash
# Compiler et lancer
go run *.go

# Ou compiler puis exécuter
go build -o power4
./power4
```

Le serveur démarre sur `http://localhost:8080`

---

## Améliorations possibles

### Niveau facile
- Ajouter un compteur de coups
- Ajouter un chronomètre
- Sauvegarder les scores dans localStorage

### Niveau intermédiaire
- Mode contre ordinateur (IA simple côté backend)
- Historique des parties (base de données)
- Thèmes personnalisables

### Niveau avancé
- Mode multijoueur en ligne avec WebSockets
- IA avancée avec algorithme minimax (backend Go)
- Replay des parties
- API REST documentée avec Swagger

---

**Bon apprentissage ! 🚀**
