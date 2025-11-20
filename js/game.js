/**
 * PUISSANCE 4 - MODULE JEU PRINCIPAL
 *
 * Ce fichier gère toute la logique de jeu côté client :
 * - Communication avec le backend Go via API REST
 * - Affichage du plateau de jeu et des jetons
 * - Gestion des animations de placement de jetons
 * - Détection et affichage de la fin de partie
 */

//#region CONFIGURATION ET CONSTANTES

/**
 * URL de base de l'API backend
 * @constant {string} URL de base pour toutes les requêtes API
 */
const API_URL = 'http://localhost:8080/api';

//#endregion

//#region VÉRIFICATION ET INITIALISATION

// Récupération de la difficulté depuis sessionStorage
const difficulty = sessionStorage.getItem('difficulty');

// Si aucune donnée n'est présente, redirection vers la page d'accueil
if (!difficulty) {
    window.location.href = '/';
}

// Application du thème visuel correspondant à la difficulté
document.body.className = difficulty;

//#endregion

//#region VARIABLES D'ÉTAT DU JEU

/**
 * Nombre de lignes du plateau de jeu
 * @type {number} Nombre de lignes du plateau
 */
let ROWS = parseInt(sessionStorage.getItem('rows'));

/**
 * Nombre de colonnes du plateau de jeu
 * @type {number} Nombre de colonnes du plateau
 */
let COLS = parseInt(sessionStorage.getItem('cols'));

/**
 * Joueur actuel ('player1' ou 'player2')
 * @type {string} Identifiant du joueur dont c'est le tour
 */
let currentPlayer = 'player1';

/**
 * Tableau 2D représentant l'état du plateau
 * Chaque case contient '', 'player1' ou 'player2'
 * @type {Array<Array<string>>} Grille de jeu bidimensionnelle
 */
let board = [];

/**
 * Indique si la partie est terminée
 * @type {boolean} État de fin de partie (true si terminée)
 */
let gameOver = false;

/**
 * Skins sélectionnés par chaque joueur
 * @type {{player1: string, player2: string}} Objet contenant les skins choisis
 */
let selectedSkins = {
    player1: sessionStorage.getItem('player1Skin'),
    player2: sessionStorage.getItem('player2Skin')
};

/**
 * Pseudos des joueurs
 * @type {{player1: string, player2: string}} Objet contenant les pseudos choisis
 */
let playerPseudos = {
    player1: sessionStorage.getItem('player1Pseudo'),
    player2: sessionStorage.getItem('player2Pseudo')
};

//#endregion

//#region COMMUNICATION AVEC LE BACKEND

/**
 * Effectue un appel à l'API backend Go
 *
 * Cette fonction gère toutes les communications avec le serveur,
 * incluant la sérialisation JSON, les headers et la gestion d'erreurs
 *
 * @param {string} endpoint - Point d'API à appeler (ex: '/game/new')
 * @param {string} [method='GET'] - Méthode HTTP à utiliser (GET, POST, etc.)
 * @param {Object|null} [data=null] - Données JSON à envoyer dans le corps de la requête (pour POST)
 * @returns {Promise<Object>} Promesse contenant la réponse JSON du serveur
 * @throws {Error} Erreur levée si la requête échoue ou si le serveur retourne une erreur
 */
async function callAPI(endpoint, method = 'GET', data = null) {
    // Configuration de base de la requête
    const options = {
        method: method,
        headers: {
            'Content-Type': 'application/json',
        },
    };

    // Ajout du body pour les requêtes POST
    if (data) {
        options.body = JSON.stringify(data);
    }

    try {
        // Envoi de la requête
        const response = await fetch(`${API_URL}${endpoint}`, options);
        const json = await response.json();

        // Vérification du statut de la réponse
        if (!response.ok) {
            throw new Error(json.error || 'Erreur serveur');
        }

        return json;
    } catch (error) {
        console.error('Erreur API:', error);
        throw error;
    }
}

//#endregion

//#region INITIALISATION DU PLATEAU

/**
 * Initialise une nouvelle partie
 *
 * Cette fonction :
 * 1. Appelle le backend pour créer une nouvelle partie
 * 2. Récupère l'état initial (avec jetons pré-remplis)
 * 3. Crée la grille visuelle HTML
 * 4. Affiche le joueur actuel
 *
 * @async
 * @returns {Promise<void>} Promesse résolue une fois l'initialisation terminée
 */
async function initBoard() {
    try {
        // Création de la partie sur le serveur
        const state = await callAPI('/game/new', 'POST', {
            rows: ROWS,
            cols: COLS,
            player1: playerPseudos.player1,
            player2: playerPseudos.player2
        });

        // Mise à jour de l'état local avec la réponse du serveur
        updateLocalState(state);

        // Création de la grille visuelle
        createBoardUI();

        // Affichage du joueur actuel
        updatePlayerDisplay();

        // Réinitialisation du message
        document.getElementById('message').textContent = '';
        document.getElementById('message').className = 'message';

    } catch (error) {
        alert('Erreur lors de la création de la partie: ' + error.message);
    }
}

/**
 * Met à jour l'état local du jeu avec les données du serveur
 *
 * Synchronise toutes les variables JavaScript avec l'état
 * retourné par le backend Go
 *
 * @param {Object} state - État du jeu retourné par le backend
 * @param {Array<Array<string>>} state.board - Grille du plateau de jeu
 * @param {string} state.currentPlayer - Identifiant du joueur dont c'est le tour
 * @param {boolean} state.gameOver - Indique si la partie est terminée
 * @param {number} state.rows - Nombre de lignes du plateau
 * @param {number} state.cols - Nombre de colonnes du plateau
 * @param {boolean} state.inverseGravity - Indique si la gravité inversée est active
 */
function updateLocalState(state) {
    board = state.board;
    currentPlayer = state.currentPlayer;
    gameOver = state.gameOver;
    ROWS = state.rows;
    COLS = state.cols;

    // Gestion de l'affichage de la gravité inversée
    if (state.inverseGravity) {
        document.body.classList.add('inverse-gravity');
    } else {
        document.body.classList.remove('inverse-gravity');
    }
}

/**
 * Crée la grille visuelle HTML du plateau de jeu
 *
 * Cette fonction génère dynamiquement toute la structure HTML
 * de la grille de jeu avec :
 * - Une div par colonne (cliquable pour déposer un jeton)
 * - Une cellule par case (affiche les jetons)
 * - Les jetons déjà présents (pré-remplis ou joués)
 */
function createBoardUI() {
    const boardElement = document.getElementById('board');
    boardElement.innerHTML = ''; // Nettoyage du plateau existant

    // Création de chaque colonne
    for (let col = 0; col < COLS; col++) {
        const columnDiv = document.createElement('div');
        columnDiv.className = 'column';

        // Ajout du gestionnaire de clic pour déposer un jeton
        columnDiv.onclick = () => dropPiece(col);

        // Création de chaque cellule de la colonne
        for (let row = 0; row < ROWS; row++) {
            const cell = document.createElement('div');
            cell.className = 'cell';

            // ID unique pour identifier la cellule (format: cell-ligne-colonne)
            cell.id = `cell-${row}-${col}`;

            // Si la case contient déjà un jeton, l'afficher
            if (board[row][col]) {
                addTokenToCell(cell, board[row][col]);
            }

            columnDiv.appendChild(cell);
        }

        boardElement.appendChild(columnDiv);
    }
}

/**
 * Ajoute visuellement un jeton dans une cellule
 *
 * Cette fonction :
 * 1. Ajoute la classe CSS du joueur (player1 ou player2)
 * 2. Crée une image avec le skin approprié
 * 3. Insère l'image dans la cellule
 *
 * @param {HTMLElement} cell - Élément HTML DOM de la cellule cible
 * @param {string} player - Identifiant du joueur possédant le jeton ('player1' ou 'player2')
 */
function addTokenToCell(cell, player) {
    // Ajout de la classe CSS pour le style
    cell.classList.add(player);

    // Création de l'image du jeton
    const img = document.createElement('img');
    if (player === 'player1') {
        img.src = `/static/tokens/${selectedSkins.player1}.png`;
    } else {
        img.src = `/static/tokens/${selectedSkins.player2}.png`;
    }

    cell.appendChild(img);
}

//#endregion

//#region LOGIQUE DE JEU

/**
 * Gère le placement d'un jeton dans une colonne
 *
 * Cette fonction est le cœur de la mécanique de jeu :
 * 1. Envoie le coup au backend
 * 2. Reçoit le nouvel état
 * 3. Anime la chute du jeton
 * 4. Vérifie la fin de partie
 *
 * @async
 * @param {number} col - Numéro de la colonne où placer le jeton (index de 0 à COLS-1)
 * @returns {Promise<void>} Promesse résolue une fois le placement et l'animation terminés
 */
async function dropPiece(col) {
    // Ne rien faire si la partie est terminée
    if (gameOver) return;

    try {
        // Envoi du coup au backend
        const state = await callAPI('/game/drop', 'POST', { col: col });

        // Vérification d'erreur (colonne pleine, etc.)
        if (state.error) {
            console.log('Coup invalide:', state.error);
            return;
        }

        // Récupération des informations sur le coup joué
        const lastMove = state.lastMove;
        const row = lastMove.row;

        // Récupération de la cellule HTML correspondante
        const cell = document.getElementById(`cell-${row}-${col}`);

        // Le backend indique quel joueur a joué
        const playerWhoMoved = state.board[row][col];

        // Ajout des classes pour le style et l'animation
        cell.classList.add(playerWhoMoved);
        cell.classList.add('dropping');

        // Ajout de l'image du jeton
        addTokenToCell(cell, state.board[row][col]);

        // Mise à jour de l'état local
        updateLocalState(state);

        // Mise à jour de l'affichage du joueur actuel
        updatePlayerDisplay();

        // Retrait de la classe d'animation après 600ms
        setTimeout(() => {
            cell.classList.remove('dropping');
        }, 600);

        // Vérification de fin de partie après l'animation
        setTimeout(() => {
            if (state.gameOver) {
                handleGameOver(state);
            }
        }, 600);

    } catch (error) {
        console.error('Erreur lors du placement du jeton:', error);
    }
}

/**
 * Gère la fin de partie (victoire ou égalité)
 *
 * Cette fonction :
 * 1. Affiche un message approprié (victoire ou égalité)
 * 2. Lance les feux d'artifice en cas de victoire
 * 3. Désactive le plateau de jeu
 *
 * @param {Object} state - État du jeu retourné par le backend
 * @param {string} state.winner - Identifiant du gagnant ('player1', 'player2' ou 'draw' en cas d'égalité)
 */
function handleGameOver(state) {
    gameOver = true;
    const message = document.getElementById('message');

    if (state.winner === 'draw') {
        // Cas d'égalité (plateau plein sans gagnant)
        message.textContent = '⚖️ Match nul ! ⚖️';
        message.className = 'message';
    } else {
        // Cas de victoire
        const winner = state.winner === 'player1' ? playerPseudos.player1 : playerPseudos.player2;
        const loser = state.winner === 'player1' ? playerPseudos.player2 : playerPseudos.player1;

        message.textContent = `Bien vue ${winner}, t'a bien soulevé le daron chauve de ${loser}`;
        message.className = 'message winner';

        // Lancement des feux d'artifice
        createFireworks();
    }
}

//#endregion

//#region AFFICHAGE ET INTERFACE

/**
 * Met à jour l'affichage du joueur actuel
 *
 * Cette fonction met à jour :
 * 1. Le pseudo du joueur affiché
 * 2. L'icône du jeton du joueur actuel
 */
function updatePlayerDisplay() {
    const playerName = document.getElementById('playerName');
    const playerIndicator = document.getElementById('playerIndicator');

    if (currentPlayer === 'player1') {
        playerName.textContent = playerPseudos.player1;
        playerIndicator.innerHTML = `<img src="/static/tokens/${selectedSkins.player1}.png">`;
    } else {
        playerName.textContent = playerPseudos.player2;
        playerIndicator.innerHTML = `<img src="/static/tokens/${selectedSkins.player2}.png">`;
    }
}

/**
 * Réinitialise complètement le jeu
 *
 * Cette fonction :
 * 1. Efface toutes les données de sessionStorage
 * 2. Redirige vers la page de sélection de difficulté
 *
 * Appelée lorsque l'utilisateur clique sur "Nouvelle Partie"
 */
function resetGame() {
    sessionStorage.clear();
    window.location.href = '/';
}

//#endregion

//#region DÉMARRAGE AUTOMATIQUE

/**
 * Initialise le jeu au chargement de la page
 * Lance automatiquement initBoard() quand le DOM est prêt
 */
window.addEventListener('DOMContentLoaded', function() {
    initBoard();
});

//#endregion
