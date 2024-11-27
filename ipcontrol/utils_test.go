package ipcontrol

func testAccConfigWithProviderIPC(config string) string {
	return serverIPC + "\n" + config
}

func testAccConfigWithProviderQIP(config string) string {
	return serverQIP + "\n" + config
}
