package savant

type mockClient struct {
	runs   [][]string
	output []string
	error  error
}

func (m *mockClient) Run(option string, args ...string) ([]string, error) {
	out := append([]string{option}, args...)
	m.runs = append(m.runs, out)
	return m.output, m.error
}

func mockSetup() *mockClient {
	scliClient = &mockClient{runs: make([][]string, 0)}
	return scliClient.(*mockClient)
}
