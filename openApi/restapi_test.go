package Openapi

import (
	"SeptimanappBackend/database"
	"SeptimanappBackend/types"
	"SeptimanappBackend/util"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// seededEvents returns the events of data/testData/horariumHelper with the ids the
// database gives them - id 1 belongs to the total-septimana event.
//
// The times are spelled out instead of being built with database.EventTime.ToTime:
// an expectation that goes through the same conversion as the code under test
// agrees with it even when the conversion is wrong.
func seededEvents() types.Events {
	berlin := util.Locale()
	return types.Events{
		{
			ID:    2,
			Start: time.Date(2019, time.July, 25, 16, 30, 0, 0, berlin),
			End:   time.Date(2019, time.July, 25, 18, 0, 0, 0, berlin),
			Names: []types.LocatedString{
				{Value: "test0", Language: "de"},
				{Value: "proba0", Language: "la"},
			},
		},
		{
			ID:    3,
			Start: time.Date(2019, time.July, 25, 20, 30, 0, 0, berlin),
			End:   time.Date(2019, time.July, 25, 21, 0, 0, 0, berlin),
			Names: []types.LocatedString{
				{Value: "test1", Language: "de"},
				{Value: "proba1", Language: "la"},
			},
		},
	}
}

func setupApi(t *testing.T) SeptimanappRestApi {
	t.Helper()
	repository, err := database.GetMockRepository()
	if err != nil {
		t.Fatalf("could not open the mock database: %v", err)
	}
	database.InitMockRepository(repository)
	return SeptimanappRestApi{repository: repository}
}

func newContext(target string) (echo.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	return echo.New().NewContext(httptest.NewRequest(http.MethodGet, target, nil), rec), rec
}

func TestSeptimanappRestApi_GetEvents(t *testing.T) {
	year2019 := 2019
	year1900 := 1900
	tests := []struct {
		name   string
		target string
		params GetEventsParams
		want   types.Events
	}{
		{"all events", "/events", GetEventsParams{}, seededEvents()},
		{"events of 2019", "/events?year=2019", GetEventsParams{Year: &year2019}, seededEvents()},
		{"no events in 1900", "/events?year=1900", GetEventsParams{Year: &year1900}, types.Events{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := setupApi(t)
			ctx, rec := newContext(tt.target)

			if !assert.NoError(t, api.GetEvents(ctx, tt.params)) {
				return
			}
			assert.Equal(t, http.StatusOK, rec.Code)

			var got types.Events
			if !assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got)) {
				return
			}
			// GetEvents collects its result from a map, so the order is undefined
			assert.True(t, util.EqualEvents(got, tt.want, true), "got %v,\n want %v", got, tt.want)
		})
	}
}

func TestSeptimanappRestApi_GetEventsId(t *testing.T) {
	// Only known ids: an unknown id currently answers 500 rather than 404, which
	// this test would freeze in place if it asserted on it.
	tests := []struct {
		name      string
		id        EventId
		wantIndex int
	}{
		{"event 2", 2, 0},
		{"event 3", 3, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := seededEvents()[tt.wantIndex]
			api := setupApi(t)
			ctx, rec := newContext(fmt.Sprintf("/events/%d", tt.id))

			if !assert.NoError(t, api.GetEventsId(ctx, tt.id)) {
				return
			}
			assert.Equal(t, http.StatusOK, rec.Code)

			var got types.Event
			if !assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got)) {
				return
			}
			assert.True(t, util.EqualEvent(&got, &want, true), "got %v,\n want %v", got, want)
		})
	}
}

// locationResponse mirrors what types.Location marshals to. It cannot be read back
// into a types.Location: MarshalJSON and UnmarshalJSON use different shapes, see
// the TODO on types.Location.UnmarshalJSON.
type locationResponse struct {
	ID              string            `json:"id"`
	OverallLocation string            `json:"overallLocation"`
	Longitude       float32           `json:"longitude"`
	Latitude        float32           `json:"latitude"`
	Altitude        float32           `json:"altitude"`
	IsMain          bool              `json:"isMain"`
	Titles          map[string]string `json:"titles"`
	Descriptions    map[string]string `json:"descriptions"`
}

func TestSeptimanappRestApi_GetLocationsId(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want locationResponse
	}{
		{"test1", "test1", locationResponse{
			ID: "test1", OverallLocation: "over1", Longitude: 100, Latitude: 1000, Altitude: 10, IsMain: true,
			Titles:       map[string]string{"de": "testTitle1de", "la": "testTitle1la"},
			Descriptions: map[string]string{"de": "testDesc1de", "la": "testDesc1la"},
		}},
		{"test2", "test2", locationResponse{
			ID: "test2", OverallLocation: "over1", Longitude: 101, Latitude: 1001, Altitude: 1, IsMain: true,
			Titles:       map[string]string{"de": "testTitle2de", "la": "testTitle2la"},
			Descriptions: map[string]string{"de": "testDesc2de", "la": "testDesc2la"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := setupApi(t)
			ctx, rec := newContext("/locations/" + tt.id)

			if !assert.NoError(t, api.GetLocationsId(ctx, tt.id)) {
				return
			}
			assert.Equal(t, http.StatusOK, rec.Code)

			var got locationResponse
			if !assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got)) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
