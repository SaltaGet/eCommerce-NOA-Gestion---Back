package jobs

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func GenerateSwagger() error {
	log.Info().Msg("Ejecutando swag init...")

	swagPath := "/home/daniel/go/bin/swag"

	cmd := exec.Command(swagPath,
		"init",
		"--generalInfo", "main.go",
		"--output", "docs",
		"--parseDependency",
		"--parseInternal",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	log.Info().Msg("Ejecutando swag fmt...")

	cmdFmt := exec.Command(swagPath, "fmt")
	cmdFmt.Stdout = os.Stdout
	cmdFmt.Stderr = os.Stderr

	if err := cmdFmt.Run(); err != nil {
		return fmt.Errorf("error al aplicar swag fmt: %w", err)
	}

	log.Info().Msg("Documentación Swagger generada y formateada correctamente.")
	return nil
}

