package pkg_test

import (
	"os"
	"testing"

	"github.com/ninadingole/gotest-ls/pkg"
	"github.com/stretchr/testify/require"
)

func Test_List(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	generateFakeFiles(t, tmpDir)

	tests := []struct {
		name       string
		fileOrDirs []string
		want       []pkg.TestDetail
		wantErr    bool
	}{
		{
			name:       "empty",
			fileOrDirs: []string{},
			want:       nil,
		},
		{
			name:       "single file",
			fileOrDirs: []string{tmpDir + "/sample/sample_test.go"},
			want: []pkg.TestDetail{
				{
					Name:         "TestSomething",
					FileName:     "sample_test.go",
					RelativePath: "sample_test.go",
					AbsolutePath: tmpDir + "/sample/sample_test.go",
					Line:         7,
					Pos:          49,
				},
			},
		},
		{
			name:       "single dir",
			fileOrDirs: []string{tmpDir + "/sample"},
			want: []pkg.TestDetail{
				{
					Name:         "TestSomething",
					FileName:     "sample_test.go",
					RelativePath: "sample/sample_test.go",
					AbsolutePath: tmpDir + "/sample/sample_test.go",
					Line:         7,
					Pos:          49,
				},
			},
		},
		{
			name:       "fail for invalid dir",
			fileOrDirs: []string{"./testdata/invalid"},
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "fail to parse invalid test file",
			fileOrDirs: []string{tmpDir + "/dummy/dummy_test.go"},
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "parse subtests correctly",
			fileOrDirs: []string{"../tests/table_test.go"},
			want:       expected,
		},
		{
			name:       "parse fuzz tests correctly",
			fileOrDirs: []string{"../tests/fuzz_test.go"},
			want: []pkg.TestDetail{
				{
					Name:         "Fuzz_Sample",
					FileName:     "fuzz_test.go",
					RelativePath: "fuzz_test.go",
					AbsolutePath: parentDir + "/tests/fuzz_test.go",
					Line:         5,
					Pos:          39,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := pkg.List(tt.fileOrDirs)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func generateFakeFiles(t *testing.T, dir string) {
	t.Helper()

	_ = os.Mkdir(dir+"/dummy", os.ModePerm)
	_ = os.Mkdir(dir+"/sample", os.ModePerm)

	//nolint:dupword, gosec
	err := os.WriteFile(dir+"/dummy/dummy_test.go", []byte(`package tests_test 

import (
	"testing"
)

dummy dummy test
`), 0o700)
	require.NoError(t, err)

	//nolint:gosec
	err = os.WriteFile(dir+"/sample/sample_test.go", []byte(`

package tests_test

import (
	"testing"
)

func TestSomething(t *testing.T) {
	t.Parallel()
	t.Skipf("Skipping...")
	t.Log("Hello, world!")
}

`), os.ModePerm)

	require.NoError(t, err)
}

var (
	pwd, _    = os.Getwd()
	parentDir = pwd[:len(pwd)-len("/pkg")]
	expected  = []pkg.TestDetail{
		{
			Name:         "Test/5_+_5_=_10",
			FileName:     "table_test.go",
			RelativePath: "table_test.go",
			AbsolutePath: parentDir + "%s/tests/table_test.go",
			Line:         23,
			Pos:          265,
		},
		{
			Name:         "Test/5_-_5_=_0",
			FileName:     "table_test.go",
			RelativePath: "table_test.go",
			AbsolutePath: parentDir + "%s/tests/table_test.go",
			Line:         30,
			Pos:          355,
		},
		{
			Name:         "Test/mixed_subtest_1",
			FileName:     "table_test.go",
			RelativePath: "table_test.go",
			AbsolutePath: parentDir + "%s/tests/table_test.go",
			Line:         12,
			Pos:          111,
		},
		{
			Name:         "Test/mixed_test_2",
			FileName:     "table_test.go",
			RelativePath: "table_test.go",
			AbsolutePath: parentDir + "%s/tests/table_test.go",
			Line:         48,
			Pos:          635,
		},
	}
)
