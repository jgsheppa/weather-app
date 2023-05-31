package weather_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jgsheppa/weather-app/weather"
)

type TestCase struct {
	want string
	got  string
}

func TestAirQualityDescription(t *testing.T) {
	t.Parallel()

	testData := []TestCase{
		{
			want: weather.GetAirQualityDescription(1),
			got:  "Good",
		},
		{
			want: weather.GetAirQualityDescription(2),
			got:  "Fair",
		},
		{
			want: weather.GetAirQualityDescription(3),
			got:  "Moderate",
		},
		{
			want: weather.GetAirQualityDescription(4),
			got:  "Poor",
		},
		{
			want: weather.GetAirQualityDescription(5),
			got:  "Very Poor",
		},
		{
			want: weather.GetAirQualityDescription(6),
			got:  "Unable to retrieve air quality data",
		},
	}

	for _, testCase := range testData {
		if !cmp.Equal(testCase.want, testCase.got) {
			t.Error(cmp.Diff(testCase.want, testCase.got))
		}
	}
}

func TestUrlBuilder(t *testing.T) {
	t.Parallel()
	params := map[string]string{
		"appid": "apiKey",
		"q":     "location",
		"limit": "5",
	}

	got := weather.BuildWeatherUrl(weather.GeocodingUrl, params)
	want := weather.BaseUrl + weather.GeocodingUrl + "?appid=apiKey&limit=5&q=location"

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
