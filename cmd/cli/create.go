package cli

import (
	"fmt"
	"log"
	"net/url"
	"os"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// longURLFlag stocke la valeur du flag --url.
var longURLFlag string

// CreateCmd représente la commande 'create'.
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Cobra gère la présence du flag avec MarkFlagRequired, mais une vérification manuelle est conservée.
		if longURLFlag == "" {
			fmt.Fprintf(os.Stderr, "Erreur: Le flag --url est requis\n")
			os.Exit(1)
		}
		// Valide le format de l'URL.
		_, err := url.ParseRequestURI(longURLFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur: URL invalide: %v\n", err)
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
		// Crée le lien court.
		link, err := linkService.CreateLink(longURLFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de la création du lien: %v\n", err)
			os.Exit(1)
		}

		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.ShortCode)
		fmt.Printf("URL courte créée avec succès:\n")
		fmt.Printf("Code: %s\n", link.ShortCode)
		fmt.Printf("URL complète: %s\n", fullShortURL)
	},
}

// init configure la commande create, ses flags, et l'ajoute à la commande racine.
func init() {
	CreateCmd.Flags().StringVar(&longURLFlag, "url", "", "URL longue à raccourcir")
	CreateCmd.MarkFlagRequired("url")
	cmd2.RootCmd.AddCommand(CreateCmd)
}
