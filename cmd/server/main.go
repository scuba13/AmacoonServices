package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/scuba13/AmacoonServices/config"

	//"github.com/scuba13/AmacoonServices/config/migrate"
	"github.com/scuba13/AmacoonServices/internal/breed"
	"github.com/scuba13/AmacoonServices/internal/cat"
	"github.com/scuba13/AmacoonServices/internal/cattery"
	"github.com/scuba13/AmacoonServices/internal/color"
	"github.com/scuba13/AmacoonServices/internal/country"
	"github.com/scuba13/AmacoonServices/internal/handler"
	"github.com/scuba13/AmacoonServices/internal/litter"
	"github.com/scuba13/AmacoonServices/internal/owner"
	"github.com/scuba13/AmacoonServices/internal/transfer"
	"github.com/scuba13/AmacoonServices/internal/federation"
	"github.com/scuba13/AmacoonServices/internal/utils"
	routes "github.com/scuba13/AmacoonServices/pkg/server"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func main() {
	// Initialize Echo and Logger
	e := echo.New()
	logger := setupLogger()
	logger.Info("Initialize Echo and Logger")
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${remote_ip} ${method} ${uri} ${status} ${error}\n",
		Output: logrus.StandardLogger().Out,
	}))

	// Load configuration data
	logger.Info("Load configuration data")
	cfg := config.LoadConfig()

	// Connect to database
	db := setupDatabase(cfg, logger)

	// Connect to Mongo
	mongo := setupMongo(cfg, logger)

	// logger.Info("Inicio Migração")
	 
	//   migrate.PopulateCountries(db,mongo)
	//   migrate.MigrateBreeds(db,mongo)
	//   migrate.MigrateColors(db,mongo)
	//   migrate.MigrateOwners(db,mongo)
	//   migrate.MigrateFederations(db,mongo)
	//   migrate.MigrateCattery(db,mongo)

	//   migrate.MigrateCats(db,mongo)
	//  migrate.MigrateCatsPattentNameToCats1(db,mongo)
	//  migrate.UpdateCatsTempWithPattensIDs2(mongo)
	// migrate.UpdateCatsWithParentIDs3(mongo)
	 

	// logger.Info("Fim Migração")
	 

	

	// Initialize repositories, handlers, and routes
	initializeApp(e, db, logger, mongo)

	// Start server
	logger.Info("Starting Server")
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}

func setupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)
	logger.Info("Logger Initialized")
	return logger
}


func setupDatabase(cfg *config.Config, logger *logrus.Logger) *gorm.DB {
	logger.Info("Connecting DB")
	db, err := config.SetupDB(cfg)
	if err != nil {
		logger.Fatalf("Failed to initialize DB connection: %v", err)
	}

	if err := db.AutoMigrate(
		&litter.LitterDB{},
		&litter.KittenDB{},
		&transfer.TransferDB{},
		&utils.ProtocolDB{},
		&utils.FilesDB{},
	); err != nil {
		logger.Fatalf("Failed to migrate database schema: %v", err)
	}
	logger.Info("Connected DB")
	return db
}

func setupMongo(cfg *config.Config, logger *logrus.Logger) *mongo.Client {
	logger.Info("Connecting Mongo")
	db, err := config.SetupMongo(cfg)
	if err != nil {
		logger.Fatalf("Failed to initialize Mongo connection: %v", err)
	}
	logger.Info("Connected Mongo")
	return db
}

func initializeApp(e *echo.Echo, db *gorm.DB, logger *logrus.Logger, mongo *mongo.Client) {
	// Initialize repositories
	logger.Info("Initialize Repositories")
	catRepo := cat.NewCatRepository(mongo,logger)
	ownerRepo := owner.NewOwnerRepository(mongo, logger)
	colorRepo := color.NewColorRepository(mongo, logger)
	litterRepo := litter.NewLitterRepository(db, logger)
	breedRepo := breed.NewBreedRepository(mongo, logger)
	countryRepo := country.NewCountryRepository(mongo, logger)
	transferepo := transfer.NewTransferRepository(db, logger)
	filesRepo := utils.NewFilesRepository(db)
	catteryRepo:= cattery.NewCatteryRepository(mongo, logger)
	federationRepo:= federation.NewFederationRepository(mongo, logger)
	logger.Info("Initialize Repositories OK")

	// Initialize converters
	logger.Info("Initialize Converters")
	litterConverter := litter.NewLitterConverter(logger)
	transferConverter := transfer.NewTransferConverter(logger)
	logger.Info("Initialize Converters OK")

	// Initialize services
	logger.Info("Initialize Services")
	litterService := litter.NewLitterService(litterRepo, filesRepo, logger, litterConverter)
	transferService := transfer.NewTransferService(transferepo, filesRepo, logger, transferConverter)
	catService := cat.NewCatService(catRepo, logger)
	breedService := breed.NewBreedService(breedRepo, logger)
	colorService := color.NewColorService(colorRepo, logger)
	countryService := country.NewCountryService(countryRepo, logger)
	ownerService := owner.NewOwnerService(ownerRepo, logger)
	catteryService := cattery.NewCatteryService(catteryRepo, logger)
	federationService:= federation.NewCatteryService(federationRepo, logger)
	logger.Info("Initialize Services OK")

	// Initialize handlers
	logger.Info("Initialize Handlers")
	catHandler := handler.NewCatHandler(catService, logger)
	ownerHandler := handler.NewOwnerHandler(ownerService, logger)
	colorHandler := handler.NewColorHandler(colorService, logger)
	litterHandler := handler.NewLitterHandler(litterService, logger)
	breedHandler := handler.NewBreedHandler(breedService, logger)
	countryHandler := handler.NewCountryHandler(countryService, logger)
	transferHandler := handler.NewTransferHandler(transferService, filesRepo, logger, transferConverter)
	catteryHandler:= handler.NewCatteryHandler(catteryService, logger)
	federationHandler:= handler.NewFederationHandler(federationService, logger)
	logger.Info("Initialize Handlers OK")

	// Initialize router and routes
	logger.Info("Initialize Router and Routes")
	routes.NewRouter(catHandler, ownerHandler, colorHandler,
		 			litterHandler, breedHandler, countryHandler,
		  			transferHandler,catteryHandler, federationHandler,
					logger, e)
	logger.Info("Initialize Router and Routes OK")

}
