package plugins

import (
	"testing"

	"github.com/jenkinsci/kubernetes-operator/pkg/log"

	"github.com/stretchr/testify/assert"
)

func TestVerifyDependencies(t *testing.T) {
	log.SetupLogger(false)

	t.Run("happy, single root plugin with one dependent plugin", func(t *testing.T) {
		basePlugins := map[Plugin][]Plugin{
			Must(New("first-root-plugin:1.0.0")): {
				Must(New("first-plugin:0.0.1")),
			},
		}
		got := VerifyDependencies(basePlugins)
		assert.Nil(t, got)
	})
	t.Run("happy, two root plugins with one depended plugin with the same version", func(t *testing.T) {
		basePlugins := map[Plugin][]Plugin{
			Must(New("first-root-plugin:1.0.0")): {
				Must(New("first-plugin:0.0.1")),
			},
			Must(New("second-root-plugin:1.0.0")): {
				Must(New("first-plugin:0.0.1")),
			},
		}
		got := VerifyDependencies(basePlugins)
		assert.Nil(t, got)
	})
	t.Run("happy, two plugin names with names with underscores", func(t *testing.T) {
		basePlugins := map[Plugin][]Plugin{
			Must(New("first_root_plugin:1.0.0")): {
				Must(New("first_plugin:0.0.1")),
			},
			Must(New("second_root_plugin:1.0.0")): {
				Must(New("first_plugin:0.0.1")),
			},
		}
		got := VerifyDependencies(basePlugins)
		assert.Nil(t, got)
	})
	t.Run("happy, two plugin names with uppercase names", func(t *testing.T) {
		basePlugins := map[Plugin][]Plugin{
			Must(New("First-Root-Plugin:1.0.0")): {
				Must(New("First_Plugin:0.0.1")),
			},
			Must(New("Second_Root_Plugin:1.0.0")): {
				Must(New("First_Plugin:0.0.1")),
			},
		}
		got := VerifyDependencies(basePlugins)
		assert.Nil(t, got)
	})
	t.Run("fail, two root plugins have different versions", func(t *testing.T) {
		basePlugins := map[Plugin][]Plugin{
			Must(New("first-root-plugin:1.0.0")): {
				Must(New("first-plugin:0.0.1")),
			},
			Must(New("first-root-plugin:2.0.0")): {
				Must(New("first-plugin:0.0.1")),
			},
		}
		got := VerifyDependencies(basePlugins)
		assert.Contains(t, got, "Plugin 'first-root-plugin:1.0.0' requires version '1.0.0' but plugin 'first-root-plugin:2.0.0' requires '2.0.0' for plugin 'first-root-plugin'")
		assert.Contains(t, got, "Plugin 'first-root-plugin:2.0.0' requires version '2.0.0' but plugin 'first-root-plugin:1.0.0' requires '1.0.0' for plugin 'first-root-plugin'")
	})
	t.Run("happy, no version collision with two separate plugins lists", func(t *testing.T) {
		basePlugins := map[Plugin][]Plugin{
			Must(New("first-root-plugin:1.0.0")): {
				Must(New("first-plugin:0.0.1")),
			},
		}
		extraPlugins := map[Plugin][]Plugin{
			Must(New("second-root-plugin:2.0.0")): {
				Must(New("first-plugin:0.0.1")),
			},
		}
		got := VerifyDependencies(basePlugins, extraPlugins)
		assert.Nil(t, got)
	})
	t.Run("fail, dependent plugins have different versions", func(t *testing.T) {
		basePlugins := map[Plugin][]Plugin{
			Must(New("first-root-plugin:1.0.0")): {
				Must(New("first-plugin:0.0.1")),
			},
			Must(New("first-root-plugin:2.0.0")): {
				Must(New("first-plugin:0.0.2")),
			},
		}
		got := VerifyDependencies(basePlugins)
		assert.Contains(t, got, "Plugin 'first-root-plugin:1.0.0' requires version '1.0.0' but plugin 'first-root-plugin:2.0.0' requires '2.0.0' for plugin 'first-root-plugin'")
		assert.Contains(t, got, "Plugin 'first-root-plugin:2.0.0' requires version '2.0.0' but plugin 'first-root-plugin:1.0.0' requires '1.0.0' for plugin 'first-root-plugin'")
		assert.Contains(t, got, "Plugin 'first-root-plugin:1.0.0' requires version '0.0.1' but plugin 'first-root-plugin:2.0.0' requires '0.0.2' for plugin 'first-plugin'")
		assert.Contains(t, got, "Plugin 'first-root-plugin:2.0.0' requires version '0.0.2' but plugin 'first-root-plugin:1.0.0' requires '0.0.1' for plugin 'first-plugin'")
	})
	t.Run("fail, root and dependent plugins have different versions", func(t *testing.T) {
		basePlugins := map[Plugin][]Plugin{
			Must(New("first-root-plugin:1.0.0")): {
				Must(New("first-plugin:0.0.1")),
			},
		}
		extraPlugins := map[Plugin][]Plugin{
			Must(New("first-root-plugin:2.0.0")): {
				Must(New("first-plugin:0.0.2")),
			},
		}
		got := VerifyDependencies(basePlugins, extraPlugins)
		assert.Contains(t, got, "Plugin 'first-root-plugin:1.0.0' requires version '1.0.0' but plugin 'first-root-plugin:2.0.0' requires '2.0.0' for plugin 'first-root-plugin'")
		assert.Contains(t, got, "Plugin 'first-root-plugin:2.0.0' requires version '2.0.0' but plugin 'first-root-plugin:1.0.0' requires '1.0.0' for plugin 'first-root-plugin'")
		assert.Contains(t, got, "Plugin 'first-root-plugin:1.0.0' requires version '0.0.1' but plugin 'first-root-plugin:2.0.0' requires '0.0.2' for plugin 'first-plugin'")
		assert.Contains(t, got, "Plugin 'first-root-plugin:2.0.0' requires version '0.0.2' but plugin 'first-root-plugin:1.0.0' requires '0.0.1' for plugin 'first-plugin'")
	})
	t.Run("happy with dash in version", func(t *testing.T) {
		basePlugins := map[Plugin][]Plugin{
			Must(New("first-root-plugin:1.0.0-1")): {
				Must(New("first-plugin:0.0.1-1")),
			},
		}
		got := VerifyDependencies(basePlugins)
		assert.Nil(t, got)
	})
}
