package global

func init() {
    // Set defaults for variables if not defined
    setEnvDefaults()

    // Validate and parse signing key
    loadSigningKey()
}
