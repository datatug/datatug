package commands

//nolint:unused // kept for backward-compatible CLI structure
//type configClientCommand struct {
//	urlConfigCommand
//}
//
////nolint:unused // legacy consoleCommand entrypoint retained for reference
//func (v *configClientCommand) Execute(_ []string) error {
//	settings, err := dtconfig.GetSettings()
//	if err != nil {
//		return fmt.Errorf("failed to get settings: %w", err)
//	}
//	if changed := v.updateUrlConfig(&settings.Client.UrlConfig); changed {
//		if err = saveConfig(settings); err != nil {
//			return fmt.Errorf("failed to save settings: %w", err)
//		}
//	}
//	return dtconfig.PrintSection(settings.Client, dtconfig.FormatYaml, os.Stdout)
//}
