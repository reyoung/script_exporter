package main

func monitor(config Config) {
	switch cfg := config.Config.(type) {
	case *CfgV1:
		monitorCfgV1(cfg)
	default:
		panic("unhandled version")
	}
}
