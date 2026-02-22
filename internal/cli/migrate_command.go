package cli

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/mph-llm-experiments/acore"
	"github.com/mph-llm-experiments/anote/internal/config"
)

// ideaMigrateCommand migrates Denote-format idea files to acore ULID format
func ideaMigrateCommand(cfg *config.Config) *Command {
	fs := flag.NewFlagSet("migrate", flag.ContinueOnError)
	applyMap := fs.String("apply-map", "", "Apply a migration map from another app")

	return &Command{
		Name:        "migrate",
		Usage:       "anote migrate [--apply-map <path>]",
		Description: "Migrate ideas from Denote format to acore format",
		Flags:       fs,
		Run: func(cmd *Command, args []string) error {
			if *applyMap != "" {
				migMap, err := acore.ReadMigrationMap(*applyMap)
				if err != nil {
					return fmt.Errorf("failed to read migration map: %w", err)
				}

				if err := acore.ApplyMappings(cfg.IdeasDirectory, migMap.Mappings); err != nil {
					return fmt.Errorf("failed to apply mappings: %w", err)
				}

				if !globalFlags.Quiet {
					fmt.Printf("Applied %d mappings from %s\n", len(migMap.Mappings), migMap.App)
				}
				return nil
			}

			// Migrate this app's files
			migMap, err := acore.MigrateDirectory(cfg.IdeasDirectory, "idea", "anote")
			if err != nil {
				return fmt.Errorf("migration failed: %w", err)
			}

			if len(migMap.Mappings) == 0 {
				if !globalFlags.Quiet {
					fmt.Println("No files to migrate.")
				}
				return nil
			}

			// Initialize the index counter from migrated files
			counter, err := acore.NewIndexCounter(cfg.IdeasDirectory, "anote")
			if err != nil {
				return fmt.Errorf("failed to create counter: %w", err)
			}
			readIndexID := func(path string) (int, error) {
				var entity struct {
					acore.Entity `yaml:",inline"`
				}
				if _, err := acore.ReadFile(path, &entity); err != nil {
					return 0, err
				}
				return entity.IndexID, nil
			}
			if err := counter.InitFromFiles("idea", readIndexID); err != nil {
				return fmt.Errorf("counter init: %w", err)
			}

			mapPath := cfg.IdeasDirectory + "/migration-map.json"
			if err := acore.WriteMigrationMap(mapPath, migMap); err != nil {
				return fmt.Errorf("failed to write migration map: %w", err)
			}

			if globalFlags.JSON {
				data, _ := json.MarshalIndent(migMap, "", "  ")
				fmt.Println(string(data))
				return nil
			}

			if !globalFlags.Quiet {
				fmt.Printf("Migrated %d ideas. Mapping saved to %s\n", len(migMap.Mappings), mapPath)
				fmt.Println("Run 'apeople migrate --apply-map " + mapPath + "' and 'atask migrate acore --apply-map " + mapPath + "' to update cross-references.")
			}
			return nil
		},
	}
}
