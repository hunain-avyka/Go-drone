package json

import (
	harness "github.com/hunain-avyka/go-spec/dist/go"
)

var JacocoJenkinsToDroneParamMapperList = []JenkinsToDroneParamMapper{
	{"changeBuildStatus", "fail_on_threshold", BoolType, nil},
	{"classPattern", "class_directories", StringType, nil},
	{"exclusionPattern", "class_exclusion_pattern", StringType, nil},
	{"inclusionPattern", "class_inclusion_pattern", StringType, nil},
	{"execPattern", "reports_path_pattern", StringType, nil},
	{"skipCopyOfSrcFiles", "skip_source_copy", BoolType, nil},
	{"sourcePattern", "source_directories", StringType, nil},
	{"sourceExclusionPattern", "source_exclusion_pattern", StringType, nil},
	{"sourceInclusionPattern", "source_inclusion_pattern", StringType, nil},
	{"minimumBranchCoverage", "threshold_branch", StringToFloat64Type, nil},
	{"minimumClassCoverage", "threshold_class", StringToFloat64Type, nil},
	{"minimumComplexityCoverage", "threshold_complexity", StringToFloat64Type, nil},
	{"minimumInstructionCoverage", "threshold_instruction", StringToFloat64Type, nil},
	{"minimumLineCoverage", "threshold_line", StringToFloat64Type, nil},
	{"minimumMethodCoverage", "threshold_method", StringToFloat64Type, nil},
	{"tool", "tool", StringType, SetJacocoTool},
	// runAlways - Missing convert logic: When: parametersMap.delegate.arguments.changeBuildStatus
	// {"runAlways", "run_always", BoolType, nil},
}

func ConvertJacoco(node Node, variables map[string]string) *harness.Step {

	step := ConvertToStepUsingParameterMapDelegate(&node, variables, JacocoJenkinsToDroneParamMapperList,
		CoverageReportImage)

	return step
}

func SetJacocoTool(node *Node, attrMap map[string]interface{}, jenkinsKey string) (interface{}, error) {
	return "jacoco", nil
}

const (
	CoverageReportImage = "plugins/coverage-report"
)
