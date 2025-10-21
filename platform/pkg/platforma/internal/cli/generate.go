package cli

import (
	"fmt"
	"hostamat/pkg/platforma/log"
	"os"
	"strings"
)

func generateCommand(args []string) {
	subject := args[0]

	switch subject {
	case "domain":
		domainName := args[1]
		log.Info("generating domain", "domain", domainName)

		err := write("./internal/"+domainName, "model.go", modelTemplate, domainName, title(domainName))
		if err != nil {
			log.Error("error", "error", err)
			return
		}

		err = write("./internal/"+domainName, "repository.go", repositoryTemplate, domainName)
		if err != nil {
			log.Error("error", "error", err)
			return
		}

		err = write("./internal/"+domainName, "service.go", serviceTemplate, domainName)
		if err != nil {
			log.Error("error", "error", err)
			return
		}

		err = write("./internal/"+domainName, "domain.go", domainTemplate, domainName)
		if err != nil {
			log.Error("error", "error", err)
			return
		}
	default:
		log.Error("can't generate subject", "subject", subject)
	}
}

func write(folder, file, template string, a ...any) error {
	// TODO: create dir only if needed
	os.MkdirAll(folder, 0755)

	data := fmt.Appendf(nil, template, a...)
	err := os.WriteFile(fmt.Sprintf("%s/%s", folder, file), data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// TODO: fix imports for platforma packages
const modelTemplate = `package %s

type %s struct{}
`

const repositoryTemplate = `package %s

import "hostamat/pkg/platforma/database"

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Schema() ([]database.Migration, database.Schema) {
	return []database.Migration{}, database.Schema{}
}
`

const serviceTemplate = `package %s

type repository interface{}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{
		repo: repo,
	}
}
`

const domainTemplate = `package %s

type Domain struct {
	Repository *Repository
	Service    *Service
}

func (d *Domain) GetRepository() any {
	return d.Repository
}

func New() *Domain {
	repository := NewRepository()
	service := NewService(repository)

	return &Domain{
		Repository: repository,
		Service:    service,
	}
}
`

func title(s string) string {
	return strings.ToUpper(string([]rune(s)[0])) + s[1:]
}
