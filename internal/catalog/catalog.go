package catalog

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

//go:embed assets/schemas assets/skills assets/coverage
var assets embed.FS

func SchemaIDs() ([]string, error) {
	return ids("assets/schemas", ".json", "")
}

func Schema(id string) ([]byte, error) {
	if !strings.HasPrefix(id, "nb/") {
		return nil, fmt.Errorf("schema id must start with nb/: %s", id)
	}
	return fs.ReadFile(assets, path.Join("assets/schemas", id+".json"))
}

func SkillIDs() ([]string, error) {
	result, err := ids("assets/skills", "SKILL.md", "")
	for i := range result {
		result[i] = strings.TrimSuffix(result[i], "/")
	}
	return result, err
}

func Skill(id string) ([]byte, error) {
	if id == "" || strings.Contains(id, "..") {
		return nil, fmt.Errorf("invalid skill id: %s", id)
	}
	return fs.ReadFile(assets, path.Join("assets/skills", id, "SKILL.md"))
}

func CoverageManifest() ([]byte, error) {
	return fs.ReadFile(assets, "assets/coverage/manifest.json")
}

func ids(root, suffix, trimPrefix string) ([]string, error) {
	var result []string
	err := fs.WalkDir(assets, root, func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(name, suffix) {
			return nil
		}
		id := strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(name, root+"/"), trimPrefix), suffix)
		result = append(result, id)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(result)
	return result, nil
}
