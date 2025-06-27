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

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// shortCodeFlag stocke la valeur du flag --code.
var shortCodeFlag string

// StatsCmd représente la commande 'stats'.
var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Affiche les statistiques (nombre de clics) pour un lien court.",
	Long: `Cette commande permet de récupérer et d'afficher le nombre total de clics
pour une URL courte spécifique en utilisant son code.

Exemple:
  url-shortener stats --code="xyz123"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Valide la présence du flag --code.
		if shortCodeFlag == "" {
			fmt.Fprintf(os.Stderr, "Erreur: Le flag --code est requis\n")
			os.Exit(1)
		}

		// Charge la configuration.
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatalf("FATAL: Configuration non chargée")
		}

		// Initialise la connexion à la base de données.
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: Échec de la connexion à la base de données: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}

		// Assure la fermeture de la connexion à la fin de l'exécution.
		defer sqlDB.Close()

		// Initialise les repositories et le service.
		linkRepo := repository.NewLinkRepository(db)
		clickRepo := repository.NewClickRepository(db)
		linkService := services.NewLinkService(linkRepo, clickRepo)

		// Récupère le lien et ses statistiques.
		link, totalClicks, err := linkService.GetLinkStats(shortCodeFlag)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Fprintf(os.Stderr, "Erreur: Aucun lien trouvé avec le code: %s\n", shortCodeFlag)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "Erreur lors de la récupération des statistiques: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Statistiques pour le code court: %s\n", link.ShortCode)
		fmt.Printf("URL longue: %s\n", link.LongURL)
		fmt.Printf("Total de clics: %d\n", totalClicks)
	},
}

// init configure la commande stats, ses flags, et l'ajoute à la commande racine.
func init() {
	StatsCmd.Flags().StringVar(&shortCodeFlag, "code", "", "Code court pour lequel récupérer les statistiques")
	StatsCmd.MarkFlagRequired("code")
	cmd2.RootCmd.AddCommand(StatsCmd)
}
