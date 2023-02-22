package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"ssstatistics/internal/config"
	d "ssstatistics/internal/domain"
	"ssstatistics/internal/service/stats"
)

var Database *gorm.DB

func CreateConnection() {
	dsn := ""

	_, exists := os.LookupEnv("POSTGRES_HOST")
	if exists {
		e := func(v string) string {
			r, _ := os.LookupEnv(v)
			return r
		}
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			e("POSTGRES_HOST"), e("POSTGRES_USER"), e("POSTGRES_PASSWORD"), e("POSTGRES_DB"), e("POSTGRES_PORT"),
		)
		fmt.Println("MY DSN ENV:", dsn)
	} else {
		c := &config.Config.DatabaseConfig
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			c.Host, c.User, c.Password, c.Dbname, c.Port,
		)
		fmt.Println("MY DSN ENV:", dsn)
	}

	if dsn == "" {
		panic("Database configuration not created!!!")
	}

	Database, _ = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			PrepareStmt:                              true,
			//Logger:                                   logger.Default.LogMode(logger.Silent),
		})

	AutoMigrate()

}

func AutoMigrate() {
	dropViews()

	err := Database.AutoMigrate(d.Models...)
	if err != nil {
		println(err.Error())
	} // Not Fucking Auto
	// Manual migrate
	for _, model := range d.Models {
		switch t := model.(type) {
		case d.MyMigrator:
			t.ColumnsMigration(Database)
		}
	}

	createViews()
}

func dropViews() {
	// Спасибо gorm
	// который пытается мигрировать каюку-то хуйню чисто по приколу
	Database.Exec(`drop materialized view if exists factions_statistics;`)
	Database.Exec(`drop materialized view if exists roles_statistics;`)
}

func createViews() {
	// Спасибо gorm
	// он не может создать вьюху, если я передаю в Exec параметр-лист. Он его куда-то отсылает
	// и постгресс ругается с "определять материализованные представления со связанными параметрами нельзя"
	// так что костыль через fmt.Sprintf()
	Database.Exec(fmt.Sprintf(`
CREATE MATERIALIZED VIEW IF NOT EXISTS factions_statistics AS
	SELECT factions.faction_name,
		   (select (select sum(a.leaves)
					from (SELECT distinct on (leave_stats.id) CASE
								  WHEN leave_time = ''
									  THEN 0
								  WHEN leave_type = 'Cryopod' AND split_part(leave_time, ':', 2)::int < 15
									  THEN 0
								  WHEN split_part(leave_time, ':', 2)::int > 5 AND split_part(leave_time, ':', 2)::int < 30
									  THEN 1
								  WHEN leave_type = 'Cryopod' AND split_part(leave_time, ':', 2)::int < 45
									  THEN 1
								  END AS leaves
						  FROM leave_stats
						  JOIN roots aaa on aaa.round_id = leave_stats.root_id
						  JOIN factions bbb on bbb.root_id = aaa.round_id
						  WHERE bbb.faction_name = factions.faction_name
							AND assigned_role IN %s
							AND leave_stats.name NOT LIKE 'maintenance drone%%') as a))::real / COUNT(factions.id)                                                          AS avg_leavers,
		   COUNT(factions.id)                                                                                                                                              AS count,
		   SUM(victory)                                                                                                                                                    AS wins,
		   SUM(victory)::real * 100 / COUNT(factions.id)::real                                                                                                                   AS winrate,
		   SUM((SELECT count(1) FROM roles where roles.owner_id = factions.id))                                                     AS members_count,
		   SUM((SELECT count(1) FROM faction_objectives fo1 where fo1.owner_id = factions.id))                                      AS total_objectives,
		   SUM((SELECT count(1) FROM faction_objectives fo1 where fo1.owner_id = factions.id and fo1.completed = 'SUCCESS'))        AS completed_objectives,
		   SUM((SELECT count(1) FROM faction_objectives fo1 where fo1.owner_id = factions.id and fo1.completed = 'SUCCESS'))::real * 100 /
				GREATEST(SUM((SELECT count(1) FROM faction_objectives fo1 where fo1.owner_id = factions.id)) ::real, 1)             AS winrate_objectives
	
			FROM factions
			group by factions.faction_name;
	`, stats.CompiledStationPositionsForSQLList))
	Database.Exec(`
	CREATE MATERIALIZED VIEW IF NOT EXISTS roles_statistics AS
		select roles.role_name,
			   COUNT(1)                                                                                                  AS count,
			   SUM(roles.victory)                                                                                        AS wins,
			   SUM(roles.victory)::real * 100 / COUNT(1)::real                                                           AS winrate,
			   SUM((SELECT count(1)
					FROM role_objectives
					WHERE roles.id = role_objectives.owner_id))                                                          AS total_objectives,
			   SUM((SELECT count(1)
					FROM role_objectives
					WHERE roles.id = role_objectives.owner_id
					  AND completed = 'SUCCESS'))                                                                        AS completed_objectives,
			   SUM((SELECT count(1)
					FROM role_objectives
					WHERE roles.id = role_objectives.owner_id
					  AND completed = 'SUCCESS'))::real * 100 /
			   GREATEST(SUM((SELECT count(1) FROM role_objectives WHERE roles.id = role_objectives.owner_id)) ::real,
						1)                                                                                               AS winrate_objectives
		
		from roles
		group by roles.role_name;
	`)
}

func RefreshMaterializedViews() []string {
	Database.Exec(`refresh materialized view factions_statistics;`)
	Database.Exec(`refresh materialized view roles_statistics;`)
	return []string{
		"Materialized views updated!",
	}
}

func Save(v any) {
	Database.Save(v)
}

func PreloadSelect(args ...string) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Select(args)
	}
}
