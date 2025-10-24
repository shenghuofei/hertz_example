package utils

import (
	"testing"
	"time"
)

func TestFormatTime_DefaultLayout(t *testing.T) {
	testTime := time.Date(2024, 6, 1, 15, 30, 45, 0, time.Local)
	expected := "2024-06-01 15:30:45"
	result := FormatTime(testTime)
	if result != expected {
		t.Errorf("FormatTime() = %v, want %v", result, expected)
	}
}

func TestFormatTime_CustomLayout(t *testing.T) {
	testTime := time.Date(2024, 6, 1, 15, 30, 45, 0, time.Local)
	expected := "2024-06-01"
	result := FormatTime(testTime, LayoutDate)
	if result != expected {
		t.Errorf("FormatTime() = %v, want %v", result, expected)
	}

	expectedTime := "15:30:45"
	resultTime := FormatTime(testTime, LayoutTime)
	if resultTime != expectedTime {
		t.Errorf("FormatTime() = %v, want %v", resultTime, expectedTime)
	}
}

func TestStringToTimestampSeconds(t *testing.T) {
	// 正确的日期时间格式
	ts, err := StringToTimestampSeconds("2024-06-01 15:30:45")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 6, 1, 15, 30, 45, 0, time.Local).Unix()
	if ts != expected {
		t.Errorf("got %d, want %d", ts, expected)
	}

	// 仅日期
	ts, err = StringToTimestampSeconds("2024-06-01")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected = time.Date(2024, 6, 1, 0, 0, 0, 0, time.Local).Unix()
	if ts != expected {
		t.Errorf("got %d, want %d", ts, expected)
	}

	// 错误格式
	_, err = StringToTimestampSeconds("invalid-date")
	if err == nil {
		t.Error("expected error for invalid date format, got nil")
	}
}
