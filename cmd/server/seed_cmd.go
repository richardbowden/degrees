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

	vcatIDs, err := seedVehicleCategories(bgCtx, queries)
	if err != nil {
		return fmt.Errorf("failed to seed vehicle categories: %w", err)
	}

	customers, err := seedDemoCustomers(bgCtx, queries, vcatIDs)
	if err != nil {
		return fmt.Errorf("failed to seed demo customers: %w", err)
	}

	completed, err := seedDemoBookings(bgCtx, queries, customers)
	if err != nil {
		return fmt.Errorf("failed to seed demo bookings: %w", err)
	}

	if err := seedDemoHistory(bgCtx, queries, completed); err != nil {
		return fmt.Errorf("failed to seed demo history: %w", err)
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

// ========================================
// Demo Data — Vehicle Categories + Price Tiers
// ========================================

func seedVehicleCategories(ctx context.Context, q *dbpg.Queries) (map[string]int64, error) {
	log.Info().Msg("seeding vehicle categories...")

	type catDef struct {
		name string
		slug string
		desc string
		sort int32
	}

	cats := []catDef{
		{"Sedan / Hatchback", "sedan-hatchback", "Small to mid-size sedans and hatchbacks", 1},
		{"SUV / Wagon", "suv-wagon", "SUVs, crossovers, and station wagons", 2},
		{"4WD / Ute", "4wd-ute", "Four-wheel drive vehicles and utility trucks", 3},
		{"Performance", "performance", "Sports cars and performance vehicles", 4},
		{"Prestige / Exotic", "prestige-exotic", "Luxury and exotic vehicles requiring specialist care", 5},
	}

	catIDs := make(map[string]int64)

	// Load any already-existing categories
	existing, err := q.ListVehicleCategories(ctx)
	if err == nil {
		for _, c := range existing {
			catIDs[c.Slug] = c.ID
		}
	}

	for _, cat := range cats {
		if _, ok := catIDs[cat.slug]; ok {
			log.Info().Str("category", cat.name).Msg("  vehicle category already exists, skipping")
			continue
		}
		created, err := q.CreateVehicleCategory(ctx, dbpg.CreateVehicleCategoryParams{
			Name:        cat.name,
			Slug:        cat.slug,
			Description: dbpg.StringToPGString(cat.desc),
			SortOrder:   cat.sort,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create vehicle category %s: %w", cat.name, err)
		}
		catIDs[cat.slug] = created.ID
		log.Info().Str("category", cat.name).Int64("id", created.ID).Msg("  vehicle category created")
	}

	if err := seedPriceTiers(ctx, q, catIDs); err != nil {
		return nil, err
	}

	return catIDs, nil
}

func seedPriceTiers(ctx context.Context, q *dbpg.Queries, catIDs map[string]int64) error {
	log.Info().Msg("seeding service price tiers...")

	// service slug -> prices per vehicle category slug (cents)
	type tierRow struct {
		svcSlug string
		prices  map[string]int64
	}

	tiers := []tierRow{
		{"basic-wash-dry", map[string]int64{
			"sedan-hatchback": 5000, "suv-wagon": 6000, "4wd-ute": 6500, "performance": 6500, "prestige-exotic": 7500,
		}},
		{"wash-protect", map[string]int64{
			"sedan-hatchback": 8000, "suv-wagon": 9500, "4wd-ute": 10000, "performance": 10000, "prestige-exotic": 12000,
		}},
		{"mini-detail", map[string]int64{
			"sedan-hatchback": 15000, "suv-wagon": 18000, "4wd-ute": 19000, "performance": 19000, "prestige-exotic": 22000,
		}},
		{"full-detail", map[string]int64{
			"sedan-hatchback": 35000, "suv-wagon": 42000, "4wd-ute": 44000, "performance": 45000, "prestige-exotic": 55000,
		}},
		{"premium-detail", map[string]int64{
			"sedan-hatchback": 50000, "suv-wagon": 60000, "4wd-ute": 63000, "performance": 65000, "prestige-exotic": 80000,
		}},
		{"stage-1-paint-correction", map[string]int64{
			"sedan-hatchback": 50000, "suv-wagon": 60000, "4wd-ute": 63000, "performance": 65000, "prestige-exotic": 75000,
		}},
		{"stage-2-paint-correction", map[string]int64{
			"sedan-hatchback": 80000, "suv-wagon": 95000, "4wd-ute": 100000, "performance": 100000, "prestige-exotic": 120000,
		}},
		{"ceramic-coating-1-year", map[string]int64{
			"sedan-hatchback": 80000, "suv-wagon": 95000, "4wd-ute": 100000, "performance": 100000, "prestige-exotic": 120000,
		}},
		{"ceramic-coating-3-year", map[string]int64{
			"sedan-hatchback": 120000, "suv-wagon": 145000, "4wd-ute": 150000, "performance": 155000, "prestige-exotic": 180000,
		}},
		{"ceramic-coating-5-year", map[string]int64{
			"sedan-hatchback": 180000, "suv-wagon": 215000, "4wd-ute": 225000, "performance": 235000, "prestige-exotic": 280000,
		}},
		{"interior-deep-clean", map[string]int64{
			"sedan-hatchback": 25000, "suv-wagon": 30000, "4wd-ute": 32000, "performance": 28000, "prestige-exotic": 35000,
		}},
		{"leather-treatment", map[string]int64{
			"sedan-hatchback": 15000, "suv-wagon": 18000, "4wd-ute": 20000, "performance": 18000, "prestige-exotic": 22000,
		}},
	}

	for _, t := range tiers {
		svc, err := q.GetServiceBySlug(ctx, dbpg.GetServiceBySlugParams{Slug: t.svcSlug})
		if err != nil {
			log.Warn().Str("service", t.svcSlug).Msg("  service not found for price tier, skipping")
			continue
		}
		for catSlug, price := range t.prices {
			catID, ok := catIDs[catSlug]
			if !ok {
				continue
			}
			if _, err := q.UpsertPriceTier(ctx, dbpg.UpsertPriceTierParams{
				ServiceID:         svc.ID,
				VehicleCategoryID: catID,
				Price:             price,
			}); err != nil {
				return fmt.Errorf("failed to upsert price tier %s/%s: %w", t.svcSlug, catSlug, err)
			}
		}
		log.Info().Str("service", t.svcSlug).Msg("  price tiers set")
	}
	return nil
}

// ========================================
// Demo Data — Customers
// ========================================

type demoCustomer struct {
	profileID int64
	vehicleID int64
	email     string
}

func seedDemoCustomers(ctx context.Context, q *dbpg.Queries, catIDs map[string]int64) ([]demoCustomer, error) {
	log.Info().Msg("seeding demo customers...")

	type custDef struct {
		email     string
		firstName string
		surname   string
		username  string
		phone     string
		address   string
		suburb    string
		postcode  string
		make      string
		model     string
		year      int32
		colour    string
		rego      string
		paintType string
		catSlug   string
		condNotes string
	}

	defs := []custDef{
		{
			email: "sarah.mitchell@example.com", firstName: "Sarah", surname: "Mitchell",
			username: "sarah.mitchell", phone: "0412 345 678",
			address: "14 Ocean Drive", suburb: "Cottesloe", postcode: "6011",
			make: "Toyota", model: "RAV4", year: 2022, colour: "Pearl White",
			rego: "1SAR234", paintType: "factory", catSlug: "suv-wagon",
			condNotes: "Minor stone chips on bonnet",
		},
		{
			email: "jake.thompson@example.com", firstName: "Jake", surname: "Thompson",
			username: "jake.thompson", phone: "0423 456 789",
			address: "22 Kings Park Road", suburb: "West Perth", postcode: "6005",
			make: "Honda", model: "Civic", year: 2021, colour: "Sonic Grey Pearl",
			rego: "1JAK567", paintType: "factory", catSlug: "sedan-hatchback",
			condNotes: "",
		},
		{
			email: "emma.nguyen@example.com", firstName: "Emma", surname: "Nguyen",
			username: "emma.nguyen", phone: "0434 567 890",
			address: "8 Hay Street", suburb: "Subiaco", postcode: "6008",
			make: "Mercedes-Benz", model: "C200", year: 2023, colour: "Obsidian Black Metallic",
			rego: "1EMM890", paintType: "factory", catSlug: "prestige-exotic",
			condNotes: "Soft paint — requires careful washing technique",
		},
		{
			email: "tom.wilson@example.com", firstName: "Tom", surname: "Wilson",
			username: "tom.wilson", phone: "0445 678 901",
			address: "55 Brighton Road", suburb: "Scarborough", postcode: "6019",
			make: "Ford", model: "Ranger XLT", year: 2022, colour: "Meteor Grey",
			rego: "1TOM123", paintType: "factory", catSlug: "4wd-ute",
			condNotes: "Used for work — mud and dust common on undercarriage",
		},
		{
			email: "lisa.chen@example.com", firstName: "Lisa", surname: "Chen",
			username: "lisa.chen", phone: "0456 789 012",
			address: "3 Market Street", suburb: "Fremantle", postcode: "6160",
			make: "Mazda", model: "CX-5", year: 2020, colour: "Machine Grey",
			rego: "1LIS456", paintType: "factory", catSlug: "suv-wagon",
			condNotes: "Has a dog — pet hair in interior",
		},
	}

	result := make([]demoCustomer, 0, len(defs))

	for _, d := range defs {
		// Get or create user
		var userID int64
		existingUser, err := q.GetUserByEmail(ctx, dbpg.GetUserByEmailParams{LoginEmail: d.email})
		if err == nil {
			userID = existingUser.ID
			log.Info().Str("email", d.email).Msg("  demo user already exists, skipping")
		} else {
			hash, err := passwordHash.HashWithDefaults("testtest", "testtest")
			if err != nil {
				return nil, fmt.Errorf("hash failed: %w", err)
			}
			user, err := q.CreateUser(ctx, dbpg.CreateUserParams{
				FirstName:      d.firstName,
				MiddleName:     dbpg.StringToPGString(""),
				Surname:        dbpg.StringToPGString(d.surname),
				Username:       d.username,
				LoginEmail:     d.email,
				PrimaryEmailID: 0,
				PasswordHash:   hash,
				SignUpStage:    "verified",
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create user %s: %w", d.email, err)
			}
			userID = user.ID
			log.Info().Int64("user_id", userID).Str("email", d.email).Msg("  demo user created")
		}

		// Get or create customer profile
		var profileID int64
		existingProfile, err := q.GetCustomerProfileByUserID(ctx, dbpg.GetCustomerProfileByUserIDParams{UserID: userID})
		if err == nil {
			profileID = existingProfile.ID
			log.Info().Int64("profile_id", profileID).Msg("  customer profile already exists, skipping")
		} else {
			profile, err := q.CreateCustomerProfile(ctx, dbpg.CreateCustomerProfileParams{
				UserID:   userID,
				Phone:    dbpg.StringToPGString(d.phone),
				Address:  dbpg.StringToPGString(d.address),
				Suburb:   dbpg.StringToPGString(d.suburb),
				Postcode: dbpg.StringToPGString(d.postcode),
				Notes:    dbpg.StringToPGString(""),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create profile %s: %w", d.email, err)
			}
			profileID = profile.ID
			log.Info().Int64("profile_id", profileID).Msg("  customer profile created")
		}

		// Get or create vehicle
		var vehicleID int64
		existingVehicles, err := q.ListVehiclesByCustomer(ctx, dbpg.ListVehiclesByCustomerParams{CustomerID: profileID})
		if err == nil && len(existingVehicles) > 0 {
			vehicleID = existingVehicles[0].ID
			log.Info().Int64("vehicle_id", vehicleID).Msg("  vehicle already exists, skipping")
		} else {
			catID := catIDs[d.catSlug]
			vehicle, err := q.CreateVehicle(ctx, dbpg.CreateVehicleParams{
				CustomerID:        profileID,
				Make:              d.make,
				Model:             d.model,
				Year:              pgtype.Int4{Int32: d.year, Valid: true},
				Colour:            dbpg.StringToPGString(d.colour),
				Rego:              dbpg.StringToPGString(d.rego),
				PaintType:         dbpg.StringToPGString(d.paintType),
				ConditionNotes:    dbpg.StringToPGString(d.condNotes),
				IsPrimary:         true,
				VehicleCategoryID: pgtype.Int8{Int64: catID, Valid: catID > 0},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create vehicle %s: %w", d.email, err)
			}
			vehicleID = vehicle.ID
			log.Info().Int64("vehicle_id", vehicleID).Msg("  vehicle created")
		}

		result = append(result, demoCustomer{profileID: profileID, vehicleID: vehicleID, email: d.email})
	}

	return result, nil
}

// ========================================
// Demo Data — Bookings
// ========================================

type demoCompletedBooking struct {
	bookingID   int64
	customerID  int64
	vehicleID   int64
	serviceSlug string
	daysAgo     int
}

func seedDemoBookings(ctx context.Context, q *dbpg.Queries, customers []demoCustomer) ([]demoCompletedBooking, error) {
	log.Info().Msg("seeding demo bookings...")

	if len(customers) < 5 {
		return nil, fmt.Errorf("expected 5 demo customers, got %d", len(customers))
	}

	now := time.Now()

	makeDate := func(daysOffset int) pgtype.Date {
		return pgtype.Date{Time: now.AddDate(0, 0, daysOffset).Truncate(24 * time.Hour), Valid: true}
	}
	makeTime := func(hour int) pgtype.Time {
		return pgtype.Time{Microseconds: int64(time.Duration(hour) * time.Hour / time.Microsecond), Valid: true}
	}

	sarah := customers[0]
	jake := customers[1]
	emma := customers[2]
	tom := customers[3]
	lisa := customers[4]

	type bookingDef struct {
		cust        demoCustomer
		svcSlug     string
		daysOffset  int
		hour        int
		status      dbpg.BookingStatus
		payStatus   dbpg.PaymentStatus
		notes       string
		daysAgo     int // >0 if completed in the past
	}

	bookings := []bookingDef{
		// Sarah — Full Detail 45 days ago (completed + fully paid)
		{sarah, "full-detail", -45, 9, dbpg.BookingStatusCompleted, dbpg.PaymentStatusFullyPaid, "Full exterior and interior detail on RAV4", 45},
		// Sarah — Stage 1 Paint Correction 20 days ago (completed + fully paid)
		{sarah, "stage-1-paint-correction", -20, 8, dbpg.BookingStatusCompleted, dbpg.PaymentStatusFullyPaid, "Light swirl removal and one-step polish", 20},
		// Jake — Basic Wash 35 days ago (completed)
		{jake, "basic-wash-dry", -35, 10, dbpg.BookingStatusCompleted, dbpg.PaymentStatusFullyPaid, "", 35},
		// Jake — Mini Detail in 21 days (confirmed, deposit paid)
		{jake, "mini-detail", 21, 9, dbpg.BookingStatusConfirmed, dbpg.PaymentStatusDepositPaid, "", 0},
		// Emma — Ceramic 3 Year 30 days ago (completed)
		{emma, "ceramic-coating-3-year", -30, 8, dbpg.BookingStatusCompleted, dbpg.PaymentStatusFullyPaid, "Full 3-year ceramic coating package including paint decontamination", 30},
		// Emma — Premium Detail in 35 days (confirmed, deposit paid)
		{emma, "premium-detail", 35, 8, dbpg.BookingStatusConfirmed, dbpg.PaymentStatusDepositPaid, "Pre-sale detail", 0},
		// Tom — Wash & Protect in 14 days (confirmed, deposit paid)
		{tom, "wash-protect", 14, 11, dbpg.BookingStatusConfirmed, dbpg.PaymentStatusDepositPaid, "", 0},
		// Tom — Full Detail in 70 days (pending payment)
		{tom, "full-detail", 70, 9, dbpg.BookingStatusPendingPayment, dbpg.PaymentStatusPending, "", 0},
		// Lisa — Interior Deep Clean in 17 days (confirmed, deposit paid)
		{lisa, "interior-deep-clean", 17, 10, dbpg.BookingStatusConfirmed, dbpg.PaymentStatusDepositPaid, "Heavy pet hair — has a golden retriever", 0},
	}

	var completed []demoCompletedBooking

	for _, b := range bookings {
		// Idempotency: skip if customer already has a booking on that date
		existing, _ := q.ListBookingsByCustomer(ctx, dbpg.ListBookingsByCustomerParams{CustomerID: b.cust.profileID})
		targetDate := now.AddDate(0, 0, b.daysOffset).Truncate(24 * time.Hour)
		alreadyExists := false
		for _, eb := range existing {
			if eb.ScheduledDate.Time.Truncate(24 * time.Hour).Equal(targetDate) {
				alreadyExists = true
				break
			}
		}
		if alreadyExists {
			log.Info().Str("service", b.svcSlug).Str("customer", b.cust.email).Msg("  booking already exists, skipping")
			// Still need to collect completed bookings for history seeding
			if b.daysAgo > 0 {
				for _, eb := range existing {
					if eb.ScheduledDate.Time.Truncate(24 * time.Hour).Equal(targetDate) {
						completed = append(completed, demoCompletedBooking{
							bookingID: eb.ID, customerID: b.cust.profileID,
							vehicleID: b.cust.vehicleID, serviceSlug: b.svcSlug, daysAgo: b.daysAgo,
						})
						break
					}
				}
			}
			continue
		}

		svc, err := q.GetServiceBySlug(ctx, dbpg.GetServiceBySlugParams{Slug: b.svcSlug})
		if err != nil {
			return nil, fmt.Errorf("service %s not found: %w", b.svcSlug, err)
		}

		subtotal := svc.BasePrice
		deposit := subtotal * 30 / 100

		booking, err := q.CreateBooking(ctx, dbpg.CreateBookingParams{
			CustomerID:            b.cust.profileID,
			VehicleID:             pgtype.Int8{Int64: b.cust.vehicleID, Valid: true},
			ScheduledDate:         makeDate(b.daysOffset),
			ScheduledTime:         makeTime(b.hour),
			EstimatedDurationMins: svc.DurationMinutes,
			Status:                b.status,
			PaymentStatus:         b.payStatus,
			Subtotal:              subtotal,
			DepositAmount:         deposit,
			TotalAmount:           subtotal,
			Notes:                 dbpg.StringToPGString(b.notes),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create booking %s for %s: %w", b.svcSlug, b.cust.email, err)
		}

		if _, err := q.CreateBookingService(ctx, dbpg.CreateBookingServiceParams{
			BookingID:      booking.ID,
			ServiceID:      svc.ID,
			PriceAtBooking: svc.BasePrice,
		}); err != nil {
			return nil, fmt.Errorf("failed to link service to booking %d: %w", booking.ID, err)
		}

		log.Info().Str("service", b.svcSlug).Str("customer", b.cust.email).
			Str("status", string(b.status)).Int64("booking_id", booking.ID).Msg("  booking created")

		if b.daysAgo > 0 {
			completed = append(completed, demoCompletedBooking{
				bookingID:   booking.ID,
				customerID:  b.cust.profileID,
				vehicleID:   b.cust.vehicleID,
				serviceSlug: b.svcSlug,
				daysAgo:     b.daysAgo,
			})
		}
	}

	return completed, nil
}

// ========================================
// Demo Data — Service History
// ========================================

func seedDemoHistory(ctx context.Context, q *dbpg.Queries, completed []demoCompletedBooking) error {
	log.Info().Msg("seeding demo service history...")

	now := time.Now()

	type historySpec struct {
		notes    []struct{ noteType, content string; visible bool }
		products []struct{ name, notes string }
		photos   []struct{ photoType, url, caption string }
	}

	// Keyed by service slug
	specs := map[string]historySpec{
		"full-detail": {
			notes: []struct{ noteType, content string; visible bool }{
				{"condition", "Vehicle arrived with moderate contamination — tree sap on roof and bonnet, light swirl marks throughout.", true},
				{"treatment", "Decontamination wash with iron remover, followed by clay bar. Two-stage machine polish on panels. Nanolicious Wash Wax applied as LSP.", true},
				{"recommendation", "Recommend a maintenance wash every 6 weeks to preserve the protection. Consider a paint correction in 12 months for the deeper scratches on driver's door.", true},
				{"follow_up", "Client asked about ceramic coating options — sent quote for 3-year package.", false},
			},
			products: []struct{ name, notes string }{
				{"Bowden's Nanolicious Wash", "Pre-wash and final wash"},
				{"Bowden's Iron Man", "Ferrous contamination removal"},
				{"Bowden's Lazy Wax", "Applied as last step product"},
				{"CarPro IronX", "Spot treatment on heavily contaminated areas"},
				{"Bowden's Clay Bar Kit", "Medium grade — full panel decontamination"},
			},
			photos: []struct{ photoType, url, caption string }{
				{"before", "https://picsum.photos/seed/fd-before-1/1200/800", "Driver side before detail — contamination and light swirls visible"},
				{"before", "https://picsum.photos/seed/fd-before-2/1200/800", "Bonnet before — tree sap and bird drop etching"},
				{"after", "https://picsum.photos/seed/fd-after-1/1200/800", "Driver side after — gloss restored, no visible swirling"},
				{"after", "https://picsum.photos/seed/fd-after-2/1200/800", "Bonnet after — sap removed, high gloss finish"},
				{"detail", "https://picsum.photos/seed/fd-detail-1/1200/800", "Interior — leather conditioned, all surfaces dressed"},
			},
		},
		"stage-1-paint-correction": {
			notes: []struct{ noteType, content string; visible bool }{
				{"condition", "Paint in good overall condition with light wash-induced swirls across all panels. No deep scratches. Bonnet has minor water spot etching.", true},
				{"treatment", "Single-stage machine polish using FLEX 3401 with medium cut pad and CarPro Reflect. Finished with Bowden's Lazy Wax on all panels.", true},
				{"recommendation", "Paint is now in excellent condition. Recommend booking a ceramic coating within 2 weeks to lock in results.", true},
			},
			products: []struct{ name, notes string }{
				{"CarPro Reflect Polish", "Single-stage correction"},
				{"Bowden's Lazy Wax", "LSP after correction"},
				{"Bowden's Wheely Clean", "Wheel cleaning"},
				{"CarPro Eraser", "Panel wipe before LSP application"},
			},
			photos: []struct{ photoType, url, caption string }{
				{"before", "https://picsum.photos/seed/s1-before-1/1200/800", "Swirl marks visible in direct sunlight"},
				{"before", "https://picsum.photos/seed/s1-before-2/1200/800", "Water spot etching on bonnet"},
				{"after", "https://picsum.photos/seed/s1-after-1/1200/800", "Panels corrected — mirror-like finish"},
				{"after", "https://picsum.photos/seed/s1-after-2/1200/800", "Bonnet after water spot removal"},
			},
		},
		"basic-wash-dry": {
			notes: []struct{ noteType, content string; visible bool }{
				{"condition", "Light dust and road grime. Interior tidy.", true},
				{"treatment", "Safe pre-wash, hand wash, and blow dry. Tyres dressed.", true},
			},
			products: []struct{ name, notes string }{
				{"Bowden's Nanolicious Wash", "Safe hand wash"},
				{"Bowden's Tyre Sheen", "Tyre dressing"},
			},
			photos: []struct{ photoType, url, caption string }{
				{"after", "https://picsum.photos/seed/bw-after-1/1200/800", "Clean and dry — ready for the road"},
			},
		},
		"ceramic-coating-3-year": {
			notes: []struct{ noteType, content string; visible bool }{
				{"condition", "Vehicle is relatively new (2023 model, 8,000km). Paint is in excellent condition with zero scratches. Factory paint protection film on bonnet.", true},
				{"treatment", "Full decontamination wash, paint inspection under LED light. Zero defects found — no correction required. Two coats CarPro Cquartz UK 3.0 applied to all painted surfaces. Boost top coat applied. Wheel arches and glass coated with CQuartz DLUX.", true},
				{"recommendation", "Avoid washing for 7 days to allow cure. Use Cquartz Reload as a maintenance spray after each wash. Book an annual inspection.", true},
				{"follow_up", "Registered coating warranty. Certificate emailed to client.", false},
			},
			products: []struct{ name, notes string }{
				{"CarPro Cquartz UK 3.0", "2 coats on all paint surfaces"},
				{"CarPro Cquartz DLUX", "Applied to glass and wheel arches"},
				{"CarPro Reload", "Boost coat / maintenance spray"},
				{"CarPro Eraser", "Panel wipe before coating"},
				{"Bowden's Nanolicious Wash", "Decontamination wash"},
				{"CarPro IronX", "Iron decontamination"},
			},
			photos: []struct{ photoType, url, caption string }{
				{"before", "https://picsum.photos/seed/cc-before-1/1200/800", "Pre-detail paint inspection — excellent condition"},
				{"after", "https://picsum.photos/seed/cc-after-1/1200/800", "Ceramic coating applied — deep gloss, water beading"},
				{"after", "https://picsum.photos/seed/cc-after-2/1200/800", "Water beading test — coating performing perfectly"},
				{"detail", "https://picsum.photos/seed/cc-detail-1/1200/800", "Glass coated — enhanced clarity and hydrophobic"},
				{"detail", "https://picsum.photos/seed/cc-detail-2/1200/800", "Wheel arches treated with DLUX"},
			},
		},
	}

	for _, b := range completed {
		// Skip if service record already exists for this booking
		existingRecords, _ := q.ListServiceRecordsByBooking(ctx, dbpg.ListServiceRecordsByBookingParams{BookingID: b.bookingID})
		if len(existingRecords) > 0 {
			log.Info().Int64("booking_id", b.bookingID).Msg("  service record already exists, skipping")
			continue
		}

		completedAt := now.AddDate(0, 0, -b.daysAgo)

		record, err := q.CreateServiceRecord(ctx, dbpg.CreateServiceRecordParams{
			BookingID:     b.bookingID,
			CustomerID:    b.customerID,
			VehicleID:     pgtype.Int8{Int64: b.vehicleID, Valid: b.vehicleID > 0},
			CompletedDate: pgtype.Timestamptz{Time: completedAt, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create service record for booking %d: %w", b.bookingID, err)
		}

		spec, ok := specs[b.serviceSlug]
		if !ok {
			log.Info().Int64("record_id", record.ID).Str("service", b.serviceSlug).Msg("  service record created (no spec)")
			continue
		}

		for _, n := range spec.notes {
			if _, err := q.CreateServiceNote(ctx, dbpg.CreateServiceNoteParams{
				ServiceRecordID:     record.ID,
				NoteType:            n.noteType,
				Content:             n.content,
				IsVisibleToCustomer: n.visible,
				CreatedBy:           pgtype.Int8{},
			}); err != nil {
				return fmt.Errorf("failed to create note for record %d: %w", record.ID, err)
			}
		}

		for _, p := range spec.products {
			if _, err := q.CreateServiceProductUsed(ctx, dbpg.CreateServiceProductUsedParams{
				ServiceRecordID: record.ID,
				ProductName:     p.name,
				Notes:           dbpg.StringToPGString(p.notes),
			}); err != nil {
				return fmt.Errorf("failed to create product for record %d: %w", record.ID, err)
			}
		}

		for _, ph := range spec.photos {
			if _, err := q.CreateServicePhoto(ctx, dbpg.CreateServicePhotoParams{
				ServiceRecordID: record.ID,
				PhotoType:       ph.photoType,
				Url:             ph.url,
				Caption:         dbpg.StringToPGString(ph.caption),
			}); err != nil {
				return fmt.Errorf("failed to create photo for record %d: %w", record.ID, err)
			}
		}

		log.Info().Int64("record_id", record.ID).Str("service", b.serviceSlug).
			Int("notes", len(spec.notes)).Int("products", len(spec.products)).Int("photos", len(spec.photos)).
			Msg("  service record created with history")
	}

	return nil
}
