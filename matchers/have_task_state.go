package matchers

import (
	"fmt"

	"github.com/cloudfoundry-incubator/receptor"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

func HaveTaskState(state string) gomega.OmegaMatcher {
	return &HaveTaskStateMatcher{
		State: state,
	}
}

type HaveTaskStateMatcher struct {
	State string
}

func (matcher *HaveTaskStateMatcher) Match(actual interface{}) (success bool, err error) {
	task, ok := actual.(receptor.TaskResponse)
	if !ok {
		return false, fmt.Errorf("HaveTaskState matcher expects a receptor.TaskResponse.  Got:\n%s", format.Object(actual, 1))
	}

	return task.State == matcher.State, nil
}

func (matcher *HaveTaskStateMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\nto have state %s", format.Object(actual, 1), matcher.State)
}

func (matcher *HaveTaskStateMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\nnot to have state %s", format.Object(actual, 1), matcher.State)
}
