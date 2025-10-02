package config

func loadTestConfig(cfg *Config) {
	cfg.Environment = "test"
	cfg.Database.Name = "mike_test_db"
}
