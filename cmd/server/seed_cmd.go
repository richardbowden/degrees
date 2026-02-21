package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/richardbowden/passwordHash"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/richardbowden/degrees/internal/dbpg"
)

const seedSchemaName = "degrees"

func seedRun(ctx *cli.Context) error {
	dbCfg := loadDBConfigFromCLI(ctx)
	dbCon, err := dbpg.NewConnection(dbCfg.ConnectionStringWithSchema(seedSchemaName), "seed")
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbCon.Close()

	queries := dbpg.New(dbCon)
	bgCtx := context.Background()

	if err := seedScheduleConfig(bgCtx, queries); err != nil {
		return fmt.Errorf("failed to seed schedule config: %w", err)
	}

	if err := seedCatalogue(bgCtx, queries); err != nil {
		return fmt.Errorf("failed to seed catalogue: %w", err)
	}

	if err := seedTestUser(bgCtx, queries); err != nil {
		return fmt.Errorf("failed to seed test user: %w", err)
	}

	log.Info().Msg("seed completed successfully")
	return nil
}

// seedScheduleConfig creates schedule config for each day of the week.
// UpdateScheduleConfig uses INSERT ON CONFLICT UPDATE, so it's idempotent.
func seedScheduleConfig(ctx context.Context, q *dbpg.Queries) error {
	log.Info().Msg("seeding schedule config...")

	dayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	for dayOfWeek := 0; dayOfWeek <= 6; dayOfWeek++ {
		isOpen := dayOfWeek >= 1 && dayOfWeek <= 6 // Mon-Sat open

		openTime := pgtype.Time{
			Microseconds: int64(7 * time.Hour / time.Microsecond), // 07:00
			Valid:        true,
		}
		closeTime := pgtype.Time{
			Microseconds: int64(17 * time.Hour / time.Microsecond), // 17:00
			Valid:        true,
		}

		_, err := q.UpdateScheduleConfig(ctx, dbpg.UpdateScheduleConfigParams{
			DayOfWeek:     int32(dayOfWeek),
			OpenTime:      openTime,
			CloseTime:     closeTime,
			IsOpen:        isOpen,
			BufferMinutes: 30,
		})
		if err != nil {
			return fmt.Errorf("failed to seed schedule for %s: %w", dayNames[dayOfWeek], err)
		}
		log.Info().Str("day", dayNames[dayOfWeek]).Bool("open", isOpen).Msg("  schedule config set")
	}

	return nil
}

func seedCatalogue(ctx context.Context, q *dbpg.Queries) error {
	log.Info().Msg("seeding catalogue...")

	type categoryDef struct {
		name string
		slug string
	}

	categories := []categoryDef{
		{"Wash & Protect", "wash-protect"},
		{"Full Detail", "full-detail"},
		{"Paint Correction", "paint-correction"},
		{"Ceramic Coating", "ceramic-coating"},
		{"Interior", "interior"},
	}

	// Create categories (skip if slug already exists)
	catIDs := make(map[string]int64)
	for i, cat := range categories {
		existing, err := q.GetCategoryBySlug(ctx, dbpg.GetCategoryBySlugParams{Slug: cat.slug})
		if err == nil {
			catIDs[cat.slug] = existing.ID
			log.Info().Str("category", cat.name).Msg("  category already exists, skipping")
			continue
		}

		created, err := q.CreateCategory(ctx, dbpg.CreateCategoryParams{
			Name:        cat.name,
			Slug:        cat.slug,
			Description: dbpg.StringToPGString(""),
			SortOrder:   int32(i + 1),
		})
		if err != nil {
			return fmt.Errorf("failed to create category %s: %w", cat.name, err)
		}
		catIDs[cat.slug] = created.ID
		log.Info().Str("category", cat.name).Int64("id", created.ID).Msg("  category created")
	}

	type serviceDef struct {
		name         string
		slug         string
		categorySlug string
		price        int64
		duration     int32
	}

	services := []serviceDef{
		{"Basic Wash & Dry", "basic-wash-dry", "wash-protect", 5000, 60},
		{"Wash & Protect", "wash-protect", "wash-protect", 8000, 90},
		{"Mini Detail", "mini-detail", "full-detail", 15000, 120},
		{"Full Detail", "full-detail", "full-detail", 35000, 240},
		{"Premium Detail", "premium-detail", "full-detail", 50000, 360},
		{"Stage 1 Paint Correction", "stage-1-paint-correction", "paint-correction", 50000, 300},
		{"Stage 2 Paint Correction", "stage-2-paint-correction", "paint-correction", 80000, 480},
		{"Ceramic Coating 1 Year", "ceramic-coating-1-year", "ceramic-coating", 80000, 360},
		{"Ceramic Coating 3 Year", "ceramic-coating-3-year", "ceramic-coating", 120000, 480},
		{"Ceramic Coating 5 Year", "ceramic-coating-5-year", "ceramic-coating", 180000, 540},
		{"Interior Deep Clean", "interior-deep-clean", "interior", 25000, 180},
		{"Leather Treatment", "leather-treatment", "interior", 15000, 120},
	}

	// Create services (skip if slug already exists)
	serviceIDs := make(map[string]int64)
	for i, svc := range services {
		catID, ok := catIDs[svc.categorySlug]
		if !ok {
			return fmt.Errorf("category %s not found for service %s", svc.categorySlug, svc.name)
		}

		existing, err := q.GetServiceBySlug(ctx, dbpg.GetServiceBySlugParams{Slug: svc.slug})
		if err == nil {
			serviceIDs[svc.slug] = existing.ID
			log.Info().Str("service", svc.name).Msg("  service already exists, skipping")
			continue
		}

		created, err := q.CreateService(ctx, dbpg.CreateServiceParams{
			CategoryID:      catID,
			Name:            svc.name,
			Slug:            svc.slug,
			Description:     dbpg.StringToPGString(svc.name + " — professional mobile detailing service."),
			ShortDesc:       dbpg.StringToPGString(svc.name),
			BasePrice:       svc.price,
			DurationMinutes: svc.duration,
			IsActive:        true,
			SortOrder:       int32(i + 1),
		})
		if err != nil {
			return fmt.Errorf("failed to create service %s: %w", svc.name, err)
		}
		serviceIDs[svc.slug] = created.ID
		log.Info().Str("service", svc.name).Int64("id", created.ID).Msg("  service created")
	}

	// Service options — attach to relevant services
	type optionDef struct {
		name     string
		price    int64
		services []string // service slugs this option applies to
	}

	options := []optionDef{
		{"Engine Bay Clean", 8000, []string{"wash-protect", "mini-detail", "full-detail", "premium-detail"}},
		{"Wheel Ceramic Coating", 15000, []string{"full-detail", "premium-detail", "ceramic-coating-1-year", "ceramic-coating-3-year", "ceramic-coating-5-year"}},
		{"Headlight Restoration", 10000, []string{"mini-detail", "full-detail", "premium-detail", "stage-1-paint-correction", "stage-2-paint-correction"}},
		{"Pet Hair Removal", 5000, []string{"basic-wash-dry", "wash-protect", "mini-detail", "full-detail", "premium-detail", "interior-deep-clean"}},
		{"Fabric Protection", 6000, []string{"interior-deep-clean", "full-detail", "premium-detail"}},
		{"Odour Removal", 8000, []string{"interior-deep-clean", "leather-treatment", "full-detail", "premium-detail"}},
	}

	for i, opt := range options {
		for _, svcSlug := range opt.services {
			svcID, ok := serviceIDs[svcSlug]
			if !ok {
				log.Warn().Str("option", opt.name).Str("service", svcSlug).Msg("  service not found for option, skipping")
				continue
			}

			// Check if option already exists on this service (by name)
			existingOpts, err := q.ListServiceOptions(ctx, dbpg.ListServiceOptionsParams{ServiceID: svcID})
			if err == nil {
				found := false
				for _, eo := range existingOpts {
					if strings.EqualFold(eo.Name, opt.name) {
						found = true
						break
					}
				}
				if found {
					continue
				}
			}

			_, err = q.CreateServiceOption(ctx, dbpg.CreateServiceOptionParams{
				ServiceID:   svcID,
				Name:        opt.name,
				Description: dbpg.StringToPGString(opt.name),
				Price:       opt.price,
				IsActive:    true,
				SortOrder:   int32(i + 1),
			})
			if err != nil {
				return fmt.Errorf("failed to create option %s on service %s: %w", opt.name, svcSlug, err)
			}
		}
		log.Info().Str("option", opt.name).Int("services", len(opt.services)).Msg("  option created")
	}

	return nil
}

func seedTestUser(ctx context.Context, q *dbpg.Queries) error {
	log.Info().Msg("seeding test user...")

	const testEmail = "test@degrees.com.au"
	const testPassword = "testtest"

	// Check if user already exists
	existingUser, err := q.GetUserByEmail(ctx, dbpg.GetUserByEmailParams{LoginEmail: testEmail})
	if err == nil {
		log.Info().Int64("user_id", existingUser.ID).Msg("  test user already exists, skipping user creation")
		return seedTestUserProfile(ctx, q, existingUser.ID)
	}

	// Hash password
	hashedPassword, err := passwordHash.HashWithDefaults(testPassword, testPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user directly
	user, err := q.CreateUser(ctx, dbpg.CreateUserParams{
		FirstName:      "Test",
		MiddleName:     dbpg.StringToPGString(""),
		Surname:        dbpg.StringToPGString("User"),
		Username:       "testuser",
		LoginEmail:     testEmail,
		PrimaryEmailID: 0,
		PasswordHash:   hashedPassword,
		SignUpStage:    "verified",
	})
	if err != nil {
		return fmt.Errorf("failed to create test user: %w", err)
	}
	log.Info().Int64("user_id", user.ID).Str("email", testEmail).Msg("  test user created")

	return seedTestUserProfile(ctx, q, user.ID)
}

func seedTestUserProfile(ctx context.Context, q *dbpg.Queries, userID int64) error {
	// Check if customer profile exists
	existingProfile, err := q.GetCustomerProfileByUserID(ctx, dbpg.GetCustomerProfileByUserIDParams{UserID: userID})
	if err == nil {
		log.Info().Int64("profile_id", existingProfile.ID).Msg("  customer profile already exists, skipping")
		return seedTestVehicle(ctx, q, existingProfile.ID)
	}

	profile, err := q.CreateCustomerProfile(ctx, dbpg.CreateCustomerProfileParams{
		UserID:   userID,
		Phone:    dbpg.StringToPGString("0400000000"),
		Address:  dbpg.StringToPGString("123 Test St"),
		Suburb:   dbpg.StringToPGString("Joondalup"),
		Postcode: dbpg.StringToPGString("6027"),
		Notes:    dbpg.StringToPGString(""),
	})
	if err != nil {
		return fmt.Errorf("failed to create customer profile: %w", err)
	}
	log.Info().Int64("profile_id", profile.ID).Msg("  customer profile created")

	return seedTestVehicle(ctx, q, profile.ID)
}

func seedTestVehicle(ctx context.Context, q *dbpg.Queries, customerID int64) error {
	// Check if vehicle already exists (by rego)
	vehicles, err := q.ListVehiclesByCustomer(ctx, dbpg.ListVehiclesByCustomerParams{CustomerID: customerID})
	if err == nil {
		for _, v := range vehicles {
			if v.Rego.String == "1TEST00" {
				log.Info().Int64("vehicle_id", v.ID).Msg("  test vehicle already exists, skipping")
				return nil
			}
		}
	}

	vehicle, err := q.CreateVehicle(ctx, dbpg.CreateVehicleParams{
		CustomerID:     customerID,
		Make:           "Jeep",
		Model:          "Wrangler",
		Year:           pgtype.Int4{Int32: 2023, Valid: true},
		Colour:         dbpg.StringToPGString("White"),
		Rego:           dbpg.StringToPGString("1TEST00"),
		PaintType:      dbpg.StringToPGString("factory"),
		ConditionNotes: dbpg.StringToPGString(""),
		IsPrimary:      true,
	})
	if err != nil {
		return fmt.Errorf("failed to create test vehicle: %w", err)
	}
	log.Info().Int64("vehicle_id", vehicle.ID).Msg("  test vehicle created")

	return nil
}
