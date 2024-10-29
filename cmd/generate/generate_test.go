package generate

import (
	"codegen/internal/flagkeys"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func TestGenerateGoSucces(t *testing.T) {
	// Constant paths.
	const memoryManifestPath = "manifest/path.json"
	const memoryOutputPath = "output/path.go"
	const packageName = "testpackage"
	const testFileManifest = "testdata/success_manifest.golden"
	const testFileGo = "testdata/success_go.golden"

	// Prepare in-memory files.
	fs := afero.NewMemMapFs()
	viper.Set(flagkeys.FileSystem, fs)
	readOsFileAndWriteToMemMap(t, testFileManifest, memoryManifestPath, fs)

	// Prepare command.
	Root.SetArgs([]string{"go",
		"--flag_manifest_path", memoryManifestPath,
		"--output_path", memoryOutputPath,
		"--package_name", packageName,
	})

	// Run command.
	Root.Execute()

	// Compare result.
	compareOutput(t, testFileGo, memoryOutputPath, fs)
}

func TestGenerateReactSucces(t *testing.T) {
	// Constant paths.
	const memoryManifestPath = "manifest/path.json"
	const memoryOutputPath = "output/path.ts"
	const testFileManifest = "testdata/success_manifest.golden"
	const testFileReact = "testdata/success_react.golden"

	// Prepare in-memory files.
	fs := afero.NewMemMapFs()
	viper.Set(flagkeys.FileSystem, fs)
	readOsFileAndWriteToMemMap(t, testFileManifest, memoryManifestPath, fs)

	// Prepare command.
	Root.SetArgs([]string{"react",
		"--flag_manifest_path", memoryManifestPath,
		"--output_path", memoryOutputPath,
	})

	// Run command.
	Root.Execute()

	// Compare result.
	compareOutput(t, testFileReact, memoryOutputPath, fs)
}

func readOsFileAndWriteToMemMap(t *testing.T, inputPath string, memPath string, memFs afero.Fs) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("error reading file %q: %v", inputPath, err)
	}
	if err := memFs.MkdirAll(filepath.Dir(memPath), 0770); err != nil {
		t.Fatalf("error creating directory %q: %v", filepath.Dir(memPath), err)
	}
	f, err := memFs.Create(memPath)
	if err != nil {
		t.Fatalf("error creating file %q: %v", memPath, err)
	}
	defer f.Close()
	writtenBytes, err := f.Write(data)
	if err != nil {
		t.Fatalf("error writing contents to file %q: %v", memPath, err)
	}
	if writtenBytes != len(data) {
		t.Fatalf("error writing entire file %v: writtenBytes != expectedWrittenBytes", memPath)
	}
}

func compareOutput(t *testing.T, testFile, memoryOutputPath string, fs afero.Fs) {
	want, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("error reading file %q: %v", testFile, err)

	}
	got, err := afero.ReadFile(fs, memoryOutputPath)
	if err != nil {
		t.Fatalf("error reading file %q: %v", memoryOutputPath, err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("output mismatch (-want +got):\n%s", diff)
	}
}
