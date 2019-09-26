package command

// Execute entry
func (m *TestAction) Execute() (interface{}, error) {
	return &TestResponse{
		Text: m.Text,
	}, nil
}
