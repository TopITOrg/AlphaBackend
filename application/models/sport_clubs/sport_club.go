package sport_clubs

import "time"

type SportClub struct {
    ID               int64      `json:"id"`
    Name             string     `json:"name"`
    Description      string     `json:"description"`
    SportType        string     `json:"sport_type"`
    TotalSlots       int32      `json:"total_slots"`
    ReservedSlots    int32      `json:"reserved_slots"`
    Location         string     `json:"location"`
    MinHoursPerWeek  int32      `json:"min_hours_per_week"`
    SkillLevel       string     `json:"skill_level"`
    CreatedBy        int64      `json:"created_by"`
    CreatedAt        time.Time  `json:"created_at"`
    UpdatedAt        time.Time  `json:"updated_at"`
    IsDeleted        bool       `json:"is_deleted"`
    Schedules        []Schedule `json:"schedules,omitempty"`
    KeyMetrics       []KeyMetric `json:"key_metrics,omitempty"`
}

type Schedule struct {
    ID        int64  `json:"id"`
    ClubID    int64  `json:"club_id"`
    DayOfWeek string `json:"day_of_week"`
    StartTime string `json:"start_time"`
    EndTime   string `json:"end_time"`
    Location  string `json:"location"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    IsDeleted bool    `json:"is_deleted"`
}

type KeyMetric struct {
    ID          int64  `json:"id"`
    ClubID      int64  `json:"club_id"`
    Name        string `json:"name"`
    Unit        string `json:"unit"`
    Description string `json:"description,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    IsDeleted   bool    `json:"is_deleted"`
}

type UpdateSportClubRequest struct {
    Name             string      `json:"name"`
    Description      string      `json:"description"`
    SportType        string      `json:"sport_type"`
    TotalSlots       int32       `json:"total_slots"`
    ReservedSlots    int32       `json:"reserved_slots"`
    Location         string      `json:"location"`
    MinHoursPerWeek  int32       `json:"min_hours_per_week"`
    SkillLevel       string      `json:"skill_level"`
    Schedules        []Schedule  `json:"schedules"`
    KeyMetrics       []KeyMetric `json:"key_metrics"`
}

type GetSportClubsResponse struct {
    Clubs  []SportClub `json:"clubs"`
    Limit  int32       `json:"limit"`
    Offset int32       `json:"offset"`
    Total  int64       `json:"total,omitempty"`
}

type Filters struct {
    SportType  string `form:"sport_type"`
    SkillLevel string `form:"skill_level"`
    Location   string `form:"location"`
    Search     string `form:"search"`
}
