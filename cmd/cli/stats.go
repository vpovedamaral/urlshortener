package cli

import (
	"errors"
	"fmt"
	"log"
	"os"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"

	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

// TODO : variable shortCodeFlag qui stockera la valeur du flag --code
var shortCodeFlag string

// StatsCmd représente la commande 'stats'
var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Affiche les statistiques (nombre de clics) pour un lien court.",
	Long: `Cette commande permet de récupérer et d'afficher le nombre total de clics
pour une URL courte spécifique en utilisant son code.

Exemple:
  url-shortener stats --code="xyz123"`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO : Valider que le flag --code a été fourni.
		if shortCodeFlag == "" {
			fmt.Fprintf(os.Stderr, "Erreur: Le flag --code est requis\n")
			os.Exit(1)
		}
		// os.Exit(1) si erreur

		// TODO : Charger la configuration chargée globalement via cmd.cfg
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatalf("FATAL: Configuration non chargée")
		}

		// TODO 3: Initialiser la connexion à la base de données SQLite avec GORM.
		// log.Fatalf si erreur

		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: Échec de la connexion à la base de données: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}

		// TODO S'assurer que la connexion est fermée à la fin de l'exécution de la commande
		defer sqlDB.Close()

		// TODO : Initialiser les repositories et services nécessaires NewLinkRepository & NewLinkService
		linkRepo := repository.NewLinkRepository(db)
		linkService := services.NewLinkService(linkRepo)

		// TODO 5: Appeler GetLinkStats pour récupérer le lien et ses statistiques.
		// Attention, la fonction retourne 3 valeurs
		// Pour l'erreur, utilisez gorm.ErrRecordNotFound
		// Si erreur, os.Exit(1)
		link, totalClicks, err := linkService.GetLinkStats(shortCodeFlag)
		if err != nil {
			// Pour l'erreur, utilisez gorm.ErrRecordNotFound
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Fprintf(os.Stderr, "Erreur: Aucun lien trouvé avec le code: %s\n", shortCodeFlag)
				os.Exit(1)
			}
			// Si erreur, os.Exit(1)
			fmt.Fprintf(os.Stderr, "Erreur lors de la récupération des statistiques: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Statistiques pour le code court: %s\n", link.ShortCode)
		fmt.Printf("URL longue: %s\n", link.LongURL)
		fmt.Printf("Total de clics: %d\n", totalClicks)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.
func init() {
	// TODO 7: Définir le flag --code pour la commande stats.
	StatsCmd.Flags().StringVar(&shortCodeFlag, "code", "", "Code court pour lequel récupérer les statistiques")

	// TODO Marquer le flag comme requis
	StatsCmd.MarkFlagRequired("code")

	// TODO : Ajouter la commande à RootCmd
	cmd2.RootCmd.AddCommand(StatsCmd)

}
