package ipcontrol

func testAccConfigWithProviderIPC(config string) string {
	return serverIPC + "\n" + config
}
