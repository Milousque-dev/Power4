/**
 * PUISSANCE 4 - MODULE FEUX D'ARTIFICE
 *
 * Ce fichier gère l'animation des feux d'artifice affichée lors de la victoire.
 * Il crée des explosions de particules colorées qui se dispersent dans toutes
 * les directions pour célébrer le gagnant.
 */

//#region FONCTION PRINCIPALE

/**
 * Crée et lance l'animation des feux d'artifice
 *
 * Cette fonction génère plusieurs explosions de particules colorées
 * sur une durée de 5 secondes pour célébrer la victoire.
 *
 * Le système fonctionne en 3 phases :
 * 1. Phase initiale : 15 feux d'artifice espacés de 250ms
 * 2. Phase continue : Nouveaux feux toutes les 400ms
 * 3. Phase d'arrêt : Fin de l'animation après 5 secondes
 *
 * Chaque feu d'artifice est composé de 50 particules qui explosent
 * dans un motif circulaire depuis un point aléatoire de l'écran.
 */
function createFireworks() {
    //#region Configuration

    // Récupération du conteneur DOM pour les feux d'artifice
    const fireworksContainer = document.getElementById('fireworks');

    /**
     * Palette de couleurs pour les particules
     * @type {string[]} Tableau de codes couleur hexadécimaux pour les particules
     */
    const colors = [
        '#ff0000',  // Rouge
        '#00ff00',  // Vert
        '#0000ff',  // Bleu
        '#ffff00',  // Jaune
        '#ff00ff',  // Magenta
        '#00ffff',  // Cyan
        '#ffa500',  // Orange
        '#ff69b4',  // Rose
        '#ffffff',  // Blanc
        '#ff1493'   // Rose foncé
    ];

    //#endregion

    //#region Création d'un feu d'artifice unique

    /**
     * Crée une seule explosion de feu d'artifice
     *
     * Cette fonction :
     * 1. Choisit une position aléatoire sur l'écran
     * 2. Sélectionne une couleur aléatoire
     * 3. Génère 50 particules disposées en cercle
     * 4. Anime chaque particule vers l'extérieur
     * 5. Supprime les particules après 1.5 secondes
     */
    function createFirework() {
        // Position aléatoire en X (largeur complète de la fenêtre)
        const x = Math.random() * window.innerWidth;

        // Position aléatoire en Y (70% supérieur de la hauteur)
        const y = Math.random() * (window.innerHeight * 0.7);

        // Sélection aléatoire d'une couleur
        const color = colors[Math.floor(Math.random() * colors.length)];

        // Création de 50 particules pour cette explosion
        for (let i = 0; i < 50; i++) {
            //#region Création de la particule

            // Création de l'élément DOM pour la particule
            const particle = document.createElement('div');
            particle.className = 'firework';

            // Positionnement au point d'origine de l'explosion
            particle.style.left = x + 'px';
            particle.style.top = y + 'px';

            // Application de la couleur
            particle.style.backgroundColor = color;
            particle.style.color = color; // Pour le box-shadow lumineux

            //#endregion

            //#region Calcul de la trajectoire

            // Calcul de l'angle de déplacement
            // Répartition circulaire uniforme des 50 particules (360 degrés)
            const angle = (Math.PI * 2 * i) / 50;

            // Vélocité aléatoire pour varier la portée de l'explosion
            // Entre 150px et 300px de distance
            const velocity = 150 + Math.random() * 150;

            // Calcul des déplacements X et Y basés sur l'angle et la vélocité
            // Utilisation de la trigonométrie pour créer un motif circulaire
            const xMove = Math.cos(angle) * velocity;
            const yMove = Math.sin(angle) * velocity;

            //#endregion

            //#region Application de l'animation

            // Définition des variables CSS personnalisées pour l'animation
            // Ces variables sont utilisées dans la keyframe CSS "explode"
            particle.style.setProperty('--x', xMove + 'px');
            particle.style.setProperty('--y', yMove + 'px');

            // Ajout de la particule au conteneur
            fireworksContainer.appendChild(particle);

            //#endregion

            //#region Nettoyage automatique

            /**
             * Supprime la particule après la fin de l'animation
             * Timeout de 1500ms pour correspondre à la durée de l'animation CSS
             */
            setTimeout(() => {
                particle.remove();
            }, 1500);

            //#endregion
        }
    }

    //#endregion

    //#region Orchestration des feux d'artifice

    /**
     * PHASE 1 : Salve initiale de feux d'artifice
     * Crée 15 explosions avec un délai progressif de 250ms entre chacune
     * Cela crée un effet de cascade au début de la célébration
     */
    for (let i = 0; i < 15; i++) {
        setTimeout(createFirework, i * 250);
    }

    /**
     * PHASE 2 : Feux d'artifice continus
     * Crée une nouvelle explosion toutes les 400ms
     * Maintient l'animation active pendant toute la durée
     */
    const interval = setInterval(createFirework, 400);

    /**
     * PHASE 3 : Arrêt de l'animation
     * Stoppe la création de nouveaux feux après 5 secondes
     * Les particules déjà créées terminent leur animation naturellement
     */
    setTimeout(() => {
        clearInterval(interval);
    }, 5000); // Durée totale de la célébration : 5 secondes

    //#endregion
}

//#endregion
