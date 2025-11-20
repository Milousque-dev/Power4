/**
 * PUISSANCE 4 - MODULE SÉLECTION DES SKINS ET PSEUDOS
 *
 * Ce fichier gère la page de sélection des skins de jetons et des pseudos
 * pour les deux joueurs. Il valide que les choix sont corrects avant de
 * permettre le démarrage de la partie.
 */

//#region VÉRIFICATION ET INITIALISATION

// Récupération de la difficulté depuis sessionStorage
const difficulty = sessionStorage.getItem('difficulty');

// Si aucune difficulté n'est sélectionnée, redirection vers la page d'accueil
if (!difficulty) {
    window.location.href = '/';
}

// Application du thème visuel correspondant à la difficulté
document.body.className = difficulty;

//#endregion

//#region VARIABLES D'ÉTAT

/**
 * Stockage des skins sélectionnés par chaque joueur
 * @type {{player1: string|null, player2: string|null}} Objet contenant les skins choisis (null si non sélectionné)
 */
let selectedSkins = {
    player1: null,  // Skin du joueur 1 (ex: "skin1", "skin2", etc.)
    player2: null   // Skin du joueur 2
};

/**
 * Stockage des pseudos saisis par chaque joueur
 * @type {{player1: string, player2: string}} Objet contenant les pseudos saisis (chaîne vide si non rempli)
 */
let playerPseudos = {
    player1: '',  // Pseudo du joueur 1
    player2: ''   // Pseudo du joueur 2
};

//#endregion

//#region GESTION DES ÉVÉNEMENTS ET VALIDATION

/**
 * Initialisation des gestionnaires d'événements au chargement de la page
 * Configure tous les écouteurs pour les skins, pseudos et le bouton de démarrage
 */
document.addEventListener('DOMContentLoaded', function() {
    //#region Récupération des éléments DOM

    const player1Options = document.querySelectorAll('#player1Skins .skin-option');
    const player2Options = document.querySelectorAll('#player2Skins .skin-option');
    const startButton = document.getElementById('startGame');
    const player1PseudoInput = document.getElementById('player1Pseudo');
    const player2PseudoInput = document.getElementById('player2Pseudo');
    const pseudoError = document.getElementById('pseudoError');

    //#endregion

    //#region Gestion des pseudos

    /**
     * Gestionnaire de saisie pour le pseudo du joueur 1
     * Met à jour l'état et vérifie si le jeu peut démarrer
     */
    player1PseudoInput.addEventListener('input', function() {
        playerPseudos.player1 = this.value.trim();
        checkIfReady();
    });

    /**
     * Gestionnaire de saisie pour le pseudo du joueur 2
     * Met à jour l'état et vérifie si le jeu peut démarrer
     */
    player2PseudoInput.addEventListener('input', function() {
        playerPseudos.player2 = this.value.trim();
        checkIfReady();
    });

    //#endregion

    //#region Gestion de la sélection des skins - Joueur 1

    /**
     * Configure les événements de clic pour tous les skins du joueur 1
     * Permet de sélectionner un skin et met à jour la disponibilité pour le joueur 2
     */
    player1Options.forEach(option => {
        option.addEventListener('click', function() {
            // Ignore le clic si le skin est désactivé (déjà pris par l'autre joueur)
            if (this.classList.contains('disabled')) return;

            // Retire la sélection précédente
            player1Options.forEach(opt => opt.classList.remove('selected'));

            // Marque ce skin comme sélectionné
            this.classList.add('selected');

            // Sauvegarde le choix
            selectedSkins.player1 = this.dataset.skin;

            // Met à jour les skins disponibles pour l'autre joueur
            updateSkinAvailability();

            // Vérifie si on peut démarrer la partie
            checkIfReady();
        });
    });

    //#endregion

    //#region Gestion de la sélection des skins - Joueur 2

    /**
     * Configure les événements de clic pour tous les skins du joueur 2
     * Permet de sélectionner un skin et met à jour la disponibilité pour le joueur 1
     */
    player2Options.forEach(option => {
        option.addEventListener('click', function() {
            // Ignore le clic si le skin est désactivé
            if (this.classList.contains('disabled')) return;

            // Retire la sélection précédente
            player2Options.forEach(opt => opt.classList.remove('selected'));

            // Marque ce skin comme sélectionné
            this.classList.add('selected');

            // Sauvegarde le choix
            selectedSkins.player2 = this.dataset.skin;

            // Met à jour les skins disponibles pour l'autre joueur
            updateSkinAvailability();

            // Vérifie si on peut démarrer la partie
            checkIfReady();
        });
    });

    //#endregion

    //#region Fonction de mise à jour de la disponibilité des skins

    /**
     * Met à jour la disponibilité des skins pour les deux joueurs
     * Désactive le skin choisi par un joueur dans la liste de l'autre joueur
     * pour éviter que les deux joueurs aient le même skin
     */
    function updateSkinAvailability() {
        // Réactive tous les skins au départ
        player1Options.forEach(opt => opt.classList.remove('disabled'));
        player2Options.forEach(opt => opt.classList.remove('disabled'));

        // Si le joueur 1 a choisi un skin, le désactiver pour le joueur 2
        if (selectedSkins.player1) {
            player2Options.forEach(opt => {
                if (opt.dataset.skin === selectedSkins.player1) {
                    opt.classList.add('disabled');
                }
            });
        }

        // Si le joueur 2 a choisi un skin, le désactiver pour le joueur 1
        if (selectedSkins.player2) {
            player1Options.forEach(opt => {
                if (opt.dataset.skin === selectedSkins.player2) {
                    opt.classList.add('disabled');
                }
            });
        }
    }

    //#endregion

    //#region Fonction de validation

    /**
     * Vérifie si toutes les conditions sont remplies pour démarrer la partie
     *
     * Conditions requises:
     * - Les deux joueurs ont choisi un skin
     * - Les deux joueurs ont entré un pseudo
     * - Les skins sont différents
     * - Les pseudos sont différents (insensible à la casse)
     *
     * Active ou désactive le bouton de démarrage en conséquence
     * Affiche un message d'erreur si les pseudos sont identiques
     */
    function checkIfReady() {
        // Cache le message d'erreur par défaut
        pseudoError.style.display = 'none';

        // Vérifie si les pseudos sont identiques (insensible à la casse)
        if (playerPseudos.player1 && playerPseudos.player2 &&
            playerPseudos.player1.toLowerCase() === playerPseudos.player2.toLowerCase()) {
            // Affiche le message d'erreur
            pseudoError.style.display = 'block';
            // Désactive le bouton
            startButton.disabled = true;
            return;
        }

        // Vérifie toutes les conditions pour activer le bouton
        if (selectedSkins.player1 && selectedSkins.player2 &&
            playerPseudos.player1 && playerPseudos.player2 &&
            selectedSkins.player1 !== selectedSkins.player2) {
            // Toutes les conditions sont remplies
            startButton.disabled = false;
        } else {
            // Au moins une condition n'est pas remplie
            startButton.disabled = true;
        }
    }

    //#endregion

    //#region Démarrage de la partie

    /**
     * Gestionnaire du bouton de démarrage de la partie
     * Sauvegarde toutes les données dans sessionStorage et redirige vers la page de jeu
     */
    startButton.addEventListener('click', function() {
        // Sauvegarde tous les choix dans sessionStorage
        sessionStorage.setItem('player1Skin', selectedSkins.player1);
        sessionStorage.setItem('player2Skin', selectedSkins.player2);
        sessionStorage.setItem('player1Pseudo', playerPseudos.player1);
        sessionStorage.setItem('player2Pseudo', playerPseudos.player2);

        // Redirection vers la page de jeu
        window.location.href = '/game';
    });

    //#endregion
});

//#endregion
