package mocks

import "github.com/so-install/internal/core/domain"

// ExecutorCall records a single call to Execute.
type ExecutorCall struct {
	Name string
	Args []string
}

// ExecutorResponse configures the return value for one Execute call.
type ExecutorResponse struct {
	Stdout string
	Stderr string
	Err    error
}

// MockExecutor records all Execute calls and returns pre-configured responses.
type MockExecutor struct {
	Calls     []ExecutorCall
	responses []ExecutorResponse
	def       ExecutorResponse
}

var _ domain.Executor = (*MockExecutor)(nil)

// AddResponse enqueues a response to be consumed by the next Execute call.
func (m *MockExecutor) AddResponse(stdout, stderr string, err error) {
	m.responses = append(m.responses, ExecutorResponse{stdout, stderr, err})
}

// SetDefault sets the response returned when the response queue is exhausted.
func (m *MockExecutor) SetDefault(stdout, stderr string, err error) {
	m.def = ExecutorResponse{stdout, stderr, err}
}

func (m *MockExecutor) Execute(name string, args ...string) (stdout, stderr string, err error) {
	m.Calls = append(m.Calls, ExecutorCall{Name: name, Args: args})
	if len(m.responses) > 0 {
		r := m.responses[0]
		m.responses = m.responses[1:]
		return r.Stdout, r.Stderr, r.Err
	}
	return m.def.Stdout, m.def.Stderr, m.def.Err
}
