package commands

//nolint:unused // kept for backward-compatible CLI structure
//type configServerCommand struct {
//	urlConfigCommand
//}
//
////nolint:unused // legacy consoleCommand entrypoint retained for reference
//func (v *configServerCommand) Execute(_ []string) error {
//	settings, err := appconfig.GetSettings()
//	if err != nil {
//		return fmt.Errorf("failed to get settings: %w", err)
//	}
//	if changed := v.updateUrlConfig(&settings.Server.UrlConfig); changed {
//		if err = saveConfig(settings); err != nil {
//			return fmt.Errorf("failed to save settings: %w", err)
//		}
//	}
//	return appconfig.PrintSection(settings.Server, appconfig.FormatYaml, os.Stdout)
//}
