package reporting

import (
	"context"
)

type MockReporter struct {
	ReportedData []ReportData
}

func (m *MockReporter) Report(_ context.Context, data ReportData) {
	m.ReportedData = append(m.ReportedData, data)
}

func NewMockReporter() *MockReporter {
	return &MockReporter{
		ReportedData: make([]ReportData, 0),
	}
}
