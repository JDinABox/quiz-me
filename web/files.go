package web

import (
	"embed"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"

	"github.com/JDinABox/quiz-me/internal/dev"
)

//go:embed dist/**
var staticFs embed.FS

var dynamicFS fs.FS

type ManifestChunk struct {
	Src            string   `json:"src,omitempty"`
	File           string   `json:"file"`
	CSS            []string `json:"css,omitempty"`
	Assets         []string `json:"assets,omitempty"`
	IsEntry        bool     `json:"isEntry,omitempty"`
	Name           string   `json:"name,omitempty"`
	Names          []string `json:"names,omitempty"`
	IsDynamicEntry bool     `json:"isDynamicEntry,omitempty"`
	Imports        []string `json:"imports,omitempty"`
	DynamicImports []string `json:"dynamicImports,omitempty"`
}

type Manifest map[string]ManifestChunk

var AssetsFs fs.FS

var ErrFileNotFound = errors.New("file not found")

var pManifest Manifest

func init() {
	if dev.IsDev {
		rootfs, err := os.OpenRoot("./web/dist")
		dynamicFS = rootfs.FS()
		if err != nil {
			slog.Error("failed to open web/dist root filesystem", "error", err)
			os.Exit(1)
		}
		if AssetsFs, err = fs.Sub(dynamicFS, "assets"); err != nil {
			slog.Error("failed to create assets sub-filesystem (dev)", "error", err)
			os.Exit(1)
		}
		return
	}

	var err error
	if AssetsFs, err = fs.Sub(staticFs, "dist/assets"); err != nil {
		slog.Error("failed to create assets sub-filesystem (prod)", "error", err)
		os.Exit(1)
	}
	f, err := staticFs.ReadFile("dist/.vite/manifest.json")
	if err != nil {
		slog.Error("failed to read manifest.json", "error", err)
		os.Exit(1)
	}
	if err = json.Unmarshal(f, &pManifest); err != nil {
		slog.Error("failed to unmarshal manifest.json", "error", err)
		os.Exit(1)
	}
}

func GetAssetUri(name string) (string, error) {
	var (
		chunk ManifestChunk
		ok    bool
	)
	if dev.IsDev {
		f, err := dynamicFS.Open(".vite/manifest.json")
		if err != nil {
			slog.Error("failed to open manifest.json in GetAssetUri", "error", err)
			os.Exit(1)
		}
		defer f.Close()

		var m Manifest

		if err = json.UnmarshalRead(f, &m); err != nil {
			slog.Error("failed to unmarshal manifest.json in GetAssetUri", "error", err)
			os.Exit(1)
		}

		chunk, ok = m[name]
	} else {
		chunk, ok = pManifest[name]
	}

	if !ok {
		return "", fmt.Errorf("[%s] error: %w", name, ErrFileNotFound)
	}
	return "/" + chunk.File, nil
}

func GetLinkPreload() (string, error) {
	styleLink, err := GetAssetUri("web/style.css")
	if err != nil {
		return "", err
	}

	scriptLink, err := GetAssetUri("web/main.ts")
	if err != nil {
		return "", err
	}

	return "<" + styleLink + ">;rel=preload;as=style;fetchpriority=high,<" + scriptLink + ">;rel=modulepreload", nil
}
