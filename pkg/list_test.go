package pkg_test

import (
	"fmt"
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
			fileOrDirs: []string{fmt.Sprintf("%s/sample/sample_test.go", tmpDir)},
			want: []pkg.TestDetail{
				{
					Name:         "TestSomething",
					FileName:     "sample_test.go",
					RelativePath: "sample_test.go",
					AbsolutePath: fmt.Sprintf("%s/sample/sample_test.go", tmpDir),
					Line:         7,
					Pos:          49,
				},
			},
		},
		{
			name:       "single dir",
			fileOrDirs: []string{fmt.Sprintf("%s/sample", tmpDir)},
			want: []pkg.TestDetail{
				{
					Name:         "TestSomething",
					FileName:     "sample_test.go",
					RelativePath: "sample/sample_test.go",
					AbsolutePath: fmt.Sprintf("%s/sample/sample_test.go", tmpDir),
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
			fileOrDirs: []string{fmt.Sprintf("%s/dummy/dummy_test.go", tmpDir)},
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "parse subtests correctly",
			fileOrDirs: []string{"../tests/table_test.go"},
			want:       expected,
		},
	}
	for _, tt := range tests {
		tt := tt

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

	_ = os.Mkdir(fmt.Sprintf("%s/dummy", dir), 0o755)
	_ = os.Mkdir(fmt.Sprintf("%s/sample", dir), 0o755)

	err := os.WriteFile(fmt.Sprintf("%s/dummy/dummy_test.go", dir), []byte(`package tests_test

import (
	"testing"
)

dummy dummy test
`), os.ModePerm)
	require.NoError(t, err)

	err = os.WriteFile(fmt.Sprintf("%s/sample/sample_test.go", dir), []byte(`package tests_test

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
		{Name: "Test/5_+_5_=_10", FileName: "table_test.go", RelativePath: "table_test.go", AbsolutePath: fmt.Sprintf("%s/tests/table_test.go", parentDir), Line: 23, Pos: 265},
		{Name: "Test/5_-_5_=_0", FileName: "table_test.go", RelativePath: "table_test.go", AbsolutePath: fmt.Sprintf("%s/tests/table_test.go", parentDir), Line: 30, Pos: 355},
		{Name: "Test/mixed_subtest_1", FileName: "table_test.go", RelativePath: "table_test.go", AbsolutePath: fmt.Sprintf("%s/tests/table_test.go", parentDir), Line: 12, Pos: 111},
		{Name: "Test/mixed_test_2", FileName: "table_test.go", RelativePath: "table_test.go", AbsolutePath: fmt.Sprintf("%s/tests/table_test.go", parentDir), Line: 48, Pos: 635},
	}
)
