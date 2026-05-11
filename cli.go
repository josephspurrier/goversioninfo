package goversioninfo

import (
	"fmt"
	"io"
	"os"
)

// CLIConfig holds all settings for generating a version info resource.
// Use NewCLIConfig to get a CLIConfig with sensible defaults.
type CLIConfig struct {
	ConfigFile          string
	OutputFile          string
	GoFile              string
	GoFilePackage       string
	PlatformSpecific    bool
	IconPath            string
	ApplicationIconPath string
	ManifestPath        string
	SkipVersionInfo     bool
	PropagateVerStrings bool

	Comment        string
	CompanyName    string
	Description    string
	FileVersion    string
	InternalName   string
	Copyright      string
	Trademark      string
	OriginalName   string
	PrivateBuild   string
	ProductName    string
	ProductVersion string
	SpecialBuild   string

	TranslationID int
	CharsetID     int

	Is64Bit bool
	IsARM   bool

	// Version override fields use -1 to mean "don't override."
	// Use NewCLIConfig to get these defaults.
	VerMajor int
	VerMinor int
	VerPatch int
	VerBuild int

	ProductVerMajor int
	ProductVerMinor int
	ProductVerPatch int
	ProductVerBuild int
}

// NewCLIConfig returns a CLIConfig with sensible defaults matching the CLI behavior.
func NewCLIConfig() CLIConfig {
	return CLIConfig{
		ConfigFile:      "versioninfo.json",
		OutputFile:      "resource.syso",
		GoFilePackage:   "main",
		VerMajor:        -1,
		VerMinor:        -1,
		VerPatch:        -1,
		VerBuild:        -1,
		ProductVerMajor: -1,
		ProductVerMinor: -1,
		ProductVerPatch: -1,
		ProductVerBuild: -1,
	}
}

// RunCLI generates version info resource files based on the provided CLIConfig.
func RunCLI(cfg CLIConfig) error {
	vi := &VersionInfo{}

	if !cfg.SkipVersionInfo {
		var input = io.ReadCloser(os.Stdin)
		if cfg.ConfigFile != "-" {
			f, err := os.Open(cfg.ConfigFile)
			if err != nil {
				return fmt.Errorf("cannot open %q: %w", cfg.ConfigFile, err)
			}
			input = f
		}

		jsonBytes, err := io.ReadAll(input)
		input.Close()
		if err != nil {
			return fmt.Errorf("error reading %q: %w", cfg.ConfigFile, err)
		}

		if err := vi.ParseJSON(jsonBytes); err != nil {
			return fmt.Errorf("could not parse the .json file: %w", err)
		}
	}

	if cfg.IconPath != "" {
		vi.IconPath = cfg.IconPath
	}
	if cfg.ApplicationIconPath != "" {
		vi.ApplicationIconPath = cfg.ApplicationIconPath
	}
	if cfg.ManifestPath != "" {
		vi.ManifestPath = cfg.ManifestPath
	}
	if cfg.Comment != "" {
		vi.StringFileInfo.Comments = cfg.Comment
	}
	if cfg.CompanyName != "" {
		vi.StringFileInfo.CompanyName = cfg.CompanyName
	}
	if cfg.Description != "" {
		vi.StringFileInfo.FileDescription = cfg.Description
	}
	if cfg.FileVersion != "" {
		vi.StringFileInfo.FileVersion = cfg.FileVersion
	}
	if cfg.InternalName != "" {
		vi.StringFileInfo.InternalName = cfg.InternalName
	}
	if cfg.Copyright != "" {
		vi.StringFileInfo.LegalCopyright = cfg.Copyright
	}
	if cfg.Trademark != "" {
		vi.StringFileInfo.LegalTrademarks = cfg.Trademark
	}
	if cfg.OriginalName != "" {
		vi.StringFileInfo.OriginalFilename = cfg.OriginalName
	}
	if cfg.PrivateBuild != "" {
		vi.StringFileInfo.PrivateBuild = cfg.PrivateBuild
	}
	if cfg.ProductName != "" {
		vi.StringFileInfo.ProductName = cfg.ProductName
	}
	if cfg.ProductVersion != "" {
		vi.StringFileInfo.ProductVersion = cfg.ProductVersion
	}
	if cfg.SpecialBuild != "" {
		vi.StringFileInfo.SpecialBuild = cfg.SpecialBuild
	}

	if cfg.TranslationID > 0 {
		vi.VarFileInfo.Translation.LangID = LangID(cfg.TranslationID)
	}
	if cfg.CharsetID > 0 {
		vi.VarFileInfo.Translation.CharsetID = CharsetID(cfg.CharsetID)
	}

	if cfg.VerMajor >= 0 {
		vi.FixedFileInfo.FileVersion.Major = cfg.VerMajor
	}
	if cfg.VerMinor >= 0 {
		vi.FixedFileInfo.FileVersion.Minor = cfg.VerMinor
	}
	if cfg.VerPatch >= 0 {
		vi.FixedFileInfo.FileVersion.Patch = cfg.VerPatch
	}
	if cfg.VerBuild >= 0 {
		vi.FixedFileInfo.FileVersion.Build = cfg.VerBuild
	}

	if cfg.ProductVerMajor >= 0 {
		vi.FixedFileInfo.ProductVersion.Major = cfg.ProductVerMajor
	}
	if cfg.ProductVerMinor >= 0 {
		vi.FixedFileInfo.ProductVersion.Minor = cfg.ProductVerMinor
	}
	if cfg.ProductVerPatch >= 0 {
		vi.FixedFileInfo.ProductVersion.Patch = cfg.ProductVerPatch
	}
	if cfg.ProductVerBuild >= 0 {
		vi.FixedFileInfo.ProductVersion.Build = cfg.ProductVerBuild
	}

	if cfg.PropagateVerStrings && vi.StringFileInfo.FileVersion != "" {
		v, err := NewFileVersion(vi.StringFileInfo.FileVersion)
		if err != nil {
			return fmt.Errorf("unexpected StringFileInfo.FileVersion format: %w", err)
		}
		vi.FixedFileInfo.FileVersion = v
	}
	if cfg.PropagateVerStrings && vi.StringFileInfo.ProductVersion != "" {
		v, err := NewFileVersion(vi.StringFileInfo.ProductVersion)
		if err != nil {
			return fmt.Errorf("unexpected StringFileInfo.ProductVersion format: %w", err)
		}
		vi.FixedFileInfo.ProductVersion = v
	}

	vi.Build()
	vi.Walk()

	if cfg.GoFile != "" {
		if err := vi.WriteGo(cfg.GoFile, cfg.GoFilePackage); err != nil {
			return fmt.Errorf("error writing Go file: %w", err)
		}
	}

	var archs []string
	if cfg.PlatformSpecific {
		archs = []string{"386", "amd64", "arm", "arm64"}
	} else {
		archs = []string{"386"}
		if cfg.IsARM {
			if cfg.Is64Bit {
				archs = []string{"arm64"}
			} else {
				archs = []string{"arm"}
			}
		} else if cfg.Is64Bit {
			archs = []string{"amd64"}
		}
	}

	for _, arch := range archs {
		fileout := cfg.OutputFile
		if len(archs) > 1 {
			fileout = fmt.Sprintf("resource_windows_%v.syso", arch)
		}
		if err := vi.WriteSyso(fileout, arch); err != nil {
			return fmt.Errorf("error writing syso: %w", err)
		}
	}

	return nil
}
