package api

import (
	log "github.com/sirupsen/logrus"

	"gitlab.com/standard-go/project/internal/app/project"
	"gitlab.com/standard-go/project/internal/app/project/service"
)



var (
	// Default Variable
	srv 				project.Service

	// Static Variable
	StoragePath 	= "./internal/storage/"
)

// SetService sets the project's data management service
func SetService(s project.Service) {
	srv = s

	switch t := s.(type) {
	case *service.Service:
		if t.Repo == nil {
			log.Fatal("service.Service's field 'Repo' should not be nil")
		}
	}
}

// Service returns the project's data management service
//
// Cara ini kurang tepat. Harusnya inisialisasi service di package driver
// (misal dari cmd/seeder atau cmd/api), lalu set-kan service itu ke sini pakai
// fungsi SetService(). Jadi fungsi ini hanya untuk dapatkan service
// yang diinisialisasi secara salah dari dalam package ini di fungsi init().
func Service() project.Service {
	return srv
}
