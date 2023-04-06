package merge_test

import (
	"fmt"
	"strings"

	. "github.com/kong/go-apiops/filebasics"
	"github.com/kong/go-apiops/merge"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Merge", func() {
	validateMerge := func(filenames []string, expected string, expectError bool) {
		GinkgoHelper()
		res, err := merge.Files(filenames)
		if err != nil {
			if expectError {
				// 'expected' is an error string to match
				Expect(err).To(MatchError(expected))
			} else {
				// 'expected' is filename of expected json output, so an error wasn't expected
				Fail(fmt.Sprintf("'%s' didn't expect error: %s", filenames, err))
			}
		} else {
			// 'expected' is filename of expected json output
			if expectError {
				// 'expected' is an error string to match, but the merge succeeded...
				Fail(fmt.Sprintf("'%s' expected error: %s", filenames, expected))
			} else {
				// 'expected' is filename of expected json output, so an error wasn't expected
				expectedResult := MustReadFile("./merge_testfiles/" + expected)
				result := MustSerialize(res, false)

				MustWriteSerializedFile("./merge_testfiles/"+
					strings.Replace(expected, "_expected.", "_generated.", -1), res, false)

				Expect(*result).To(MatchJSON(*expectedResult))
			}
		}
	}

	Describe("merges files", func() {
		It("attempt 1: original order", func() {
			// This tests the order of the resulting file, but also the version of the
			// final file
			fileList := []string{
				"./merge_testfiles/file1.yml",
				"./merge_testfiles/file2.yml",
				"./merge_testfiles/file3.yml",
			}
			expected := "test1_expected.json"
			expectErr := false

			validateMerge(fileList, expected, expectErr)
		})

		It("attempt 2: same files, different order", func() {
			// This tests the order of the resulting file (different from attempt 1),
			// but also the version of the final file (same as attempt 1)
			fileList := []string{
				"./merge_testfiles/file3.yml",
				"./merge_testfiles/file2.yml",
				"./merge_testfiles/file1.yml",
			}
			expected := "test2_expected.json"
			expectErr := false

			validateMerge(fileList, expected, expectErr)
		})

		It("with incompatible versions errors", func() {
			fileList := []string{
				"./merge_testfiles/file1.yml",
				"./merge_testfiles/badversion.yml",
			}
			expected := "failed to merge ./merge_testfiles/badversion.yml: files are incompatible; " +
				"major versions are incompatible; 3.0 and 1.0"
			expectErr := true

			validateMerge(fileList, expected, expectErr)
		})

		It("with incompatible '_transform' errors", func() {
			fileList := []string{
				"./merge_testfiles/file1.yml",
				"./merge_testfiles/transform_false.yml",
			}
			expected := "failed to merge ./merge_testfiles/transform_false.yml: files are incompatible; " +
				"files with '_transform: true' (default) and '_transform: false' are not compatible"
			expectErr := true

			validateMerge(fileList, expected, expectErr)
		})
	})
})