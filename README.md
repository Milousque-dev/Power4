# 🎮 Puissance 4

Un jeu de Puissance 4 interactif avec backend Go et frontend HTML/CSS/JavaScript.

## 📋 Fonctionnalités

### 🎯 Modes de difficulté
- **Facile** : Grille 6×7 avec thème orange
- **Normal** : Grille 7×8 avec thème cyber vert
- **Difficile** : Grille 8×9 avec thème anime rose/cyan

### 🎨 Personnalisation
- 8 skins de jetons différents
- Pseudos personnalisables
- Les deux joueurs ne peuvent pas choisir le même jeton

### ✨ Animations
- Animation de chute réaliste des jetons
- Feux d'artifice spectaculaires lors de la victoire

## 📁 Structure du projet

```
Power4/
├── BACKEND (Go)
│   ├── main.go           # Serveur HTTP et routes API
│   ├── game.go           # Logique du jeu
│   ├── game_manager.go   # Gestionnaire d'état
│   └── go.mod            # Module Go
│
├── FRONTEND
│   ├── index.html        # Structure HTML
│   ├── css/styles.css    # Styles et animations
│   ├── js/
│   │   ├── game.js       # Communication API
│   │   ├── ui.js         # Interface utilisateur
│   │   └── fireworks.js  # Animation victoire
│   └── static/
│       ├── maps/         # Fonds par difficulté
│       └── tokens/       # Skins des jetons
│
└── DOCUMENTATION
    ├── README.md         # Ce fichier
    └── TECHNICAL.md      # Documentation technique
```

## 🚀 Installation et utilisation

### Prérequis
- **Go 1.21+** ([Télécharger Go](https://go.dev/dl/))
- Navigateur web moderne

### Lancer le jeu

```bash
# 1. Démarrer le serveur backend
go run *.go

# 2. Ouvrir le navigateur à l'adresse
http://localhost:8080
```

### Jouer
1. Choisissez votre difficulté
2. Entrez les pseudos et sélectionnez les jetons
3. Cliquez sur une colonne pour jouer !

## 💻 Architecture

### Backend Go - Logique du jeu

**game.go** : Logique complète du Puissance 4
- Structure `Game` avec plateau 2D
- `DropPiece()` : Place un jeton et valide le coup
- `checkWin()` : Détecte les 4 alignés (horizontal, vertical, diagonales)
- `checkDraw()` : Vérifie si le plateau est plein

**game_manager.go** : Gestion de l'état
- Gère la partie en cours
- Handlers HTTP pour les requêtes API

**main.go** : Serveur HTTP
- Écoute sur le port 8080
- Routes API : `/api/game/new`, `/api/game/drop`, `/api/game/state`, `/api/game/reset`
- Sert les fichiers statiques (HTML, CSS, JS)

### Frontend JavaScript - Interface

**game.js** : Communication avec le backend
- Appels API via `fetch()`
- Mise à jour de l'interface avec les réponses
- Animations de chute des jetons
- **Aucune logique de jeu** (tout est dans le backend Go)

**ui.js** : Gestion de l'interface
- Sélection difficulté et configuration grille
- Gestion skins/pseudos
- Affichage du joueur actuel

**fireworks.js** : Animation de victoire
- Particules explosives avec couleurs aléatoires
- Animation continue pendant 5 secondes

## 🎮 API REST

| Méthode | Endpoint | Body | Description |
|---------|----------|------|-------------|
| POST | `/api/game/new` | `{rows, cols, player1, player2}` | Créer une partie |
| POST | `/api/game/drop` | `{col}` | Jouer un coup |
| GET | `/api/game/state` | - | Obtenir l'état actuel |
| POST | `/api/game/reset` | - | Réinitialiser |

**Exemple de réponse** :
```json
{
  "rows": 6,
  "cols": 7,
  "board": [["", "", ...], ...],
  "currentPlayer": "player1",
  "player1": "Alice",
  "player2": "Bob",
  "gameOver": false,
  "winner": "",
  "lastMove": {"row": 5, "col": 3}
}
```

## 🔧 Personnalisation

### Ajouter un skin de jeton
1. Ajoutez votre image PNG dans `static/tokens/`
2. Nommez-la `skinX.png`
3. Ajoutez dans `index.html` :
```html
<div class="skin-option" data-skin="skinX">
    <img src="static/tokens/skinX.png" alt="Skin X">
</div>
```

### Modifier les dimensions de grille
Dans `ui.js`, fonction `selectDifficulty()` :
```javascript
case 'easy':
    ROWS = 6;  // Lignes
    COLS = 7;  // Colonnes
```

### Changer les couleurs
Dans `css/styles.css`, sections `body.easy`, `body.normal`, `body.hard`

## 🐛 Résolution de problèmes

**Les images ne s'affichent pas**
- Vérifiez que le dossier `static/` existe

**Le serveur ne démarre pas**
- Vérifiez que le port 8080 est libre
- Vérifiez que Go est installé : `go version`

**Le jeu ne répond pas**
- Ouvrez la console du navigateur (F12)
- Vérifiez que le backend Go est lancé

## 📝 Contraintes du projet

Ce projet est un exercice académique avec la contrainte suivante :
- ❌ **Interdiction** d'utiliser JavaScript pour la logique du jeu
- ✅ **Obligation** d'utiliser Golang pour toute la logique métier
- ✅ JavaScript uniquement pour l'interface utilisateur (UI/UX)

**Résultat** : Séparation stricte Backend (Go) / Frontend (JS) via API REST.

---

**Amusez-vous bien ! 🎮**
