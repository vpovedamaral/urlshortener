package cli

import (
	"fmt"
	"log"
	"net/url" // Pour valider le format de l'URL
	"os"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

// TODO : Faire une variable longURLFlag qui stockera la valeur du flag --url
var longURLFlag string

// CreateCmd représente la commande 'create'
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO 1: Valider que le flag --url a été fourni.
		if longURLFlag == "" {
			fmt.Fprintf(os.Stderr, "Erreur: Le flag --url est requis\n")
			os.Exit(1)
		}
		// TODO Validation basique du format de l'URL avec le package url et la fonction ParseRequestURI
		// si erreur, os.Exit(1)
		_, err := url.ParseRequestURI(longURLFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur: URL invalide: %v\n", err)
			os.Exit(1)
		}
		// TODO : Charger la configuration chargée globalement via cmd.cfg
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatalf("FATAL: Configuration non chargée")
		}

		// TODO : Initialiser la connexion à la base de données SQLite.
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
		clickRepo := repository.NewClickRepository(db)
		linkService := services.NewLinkService(linkRepo, clickRepo)
		// TODO : Appeler le LinkService et la fonction CreateLink pour créer le lien court.
		link, err := linkService.CreateLink(longURLFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de la création du lien: %v\n", err)
			os.Exit(1)
		}
		// os.Exit(1) si erreur

		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.ShortCode)
		fmt.Printf("URL courte créée avec succès:\n")
		fmt.Printf("Code: %s\n", link.ShortCode)
		fmt.Printf("URL complète: %s\n", fullShortURL)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.
func init() {
	// TODO : Définir le flag --url pour la commande create.
	CreateCmd.Flags().StringVar(&longURLFlag, "url", "", "URL longue à raccourcir")
	// TODO :  Marquer le flag comme requis
	CreateCmd.MarkFlagRequired("url")
	// TODO : Ajouter la commande à RootCmd
	cmd2.RootCmd.AddCommand(CreateCmd)
}
