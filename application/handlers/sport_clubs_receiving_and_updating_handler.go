package handlers

import (
    "fmt"
    "net/http"
    "strconv"
    "sport_platform/application/models/claims"
    "sport_platform/application/models/sport_clubs"
    "sport_platform/internal/mapper"
    "sport_platform/internal/middleware"
    "sport_platform/internal/service_wrapper"

    "github.com/gin-gonic/gin"
)

func GetSportClubsHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
    // Параметры пагинации и фильтров
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
    offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
    
    var filters sport_clubs.Filters
    if err := ctx.ShouldBindQuery(&filters); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Invalid query parameters",
        })
        return
    }

    // Получаем клубы с фильтрами
    clubs, err := wrapper.Db.Queries.GetSportClubs(ctx, map[string]interface{}{
        "limit":       int32(limit),
        "offset":      int32(offset),
        "sport_type":  filters.SportType,
        "skill_level": filters.SkillLevel,
        "location":    filters.Location,
        "search":      filters.Search,
    })
    if err != nil {
        fmt.Printf("Error fetching sport clubs: %s\n", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to fetch sport clubs",
        })
        return
    }

    // Преобразуем в модель ответа
    var response sport_clubs.GetSportClubsResponse
    response.Limit = int32(limit)
    response.Offset = int32(offset)
    response.Clubs = make([]sport_clubs.SportClub, 0, len(clubs))

    for _, club := range clubs {
        var sportClub sport_clubs.SportClub
        if err := mapper.Mapper{}.Map(&sportClub, club); err != nil {
            fmt.Printf("Error mapping club: %s\n", err)
            continue
        }

        // Загружаем расписание
        schedules, err := wrapper.Db.Queries.GetClubSchedules(ctx, club.ID)
        if err == nil {
            for _, schedule := range schedules {
                sportClub.Schedules = append(sportClub.Schedules, sport_clubs.Schedule{
                    ID:        schedule.ID,
                    ClubID:    schedule.ClubID,
                    DayOfWeek: schedule.DayOfWeek,
                    StartTime: schedule.StartTime,
                    EndTime:   schedule.EndTime,
                    Location:  schedule.Location,
                    CreatedAt: schedule.CreatedAt,
                    UpdatedAt: schedule.UpdatedAt,
                    IsDeleted: schedule.IsDeleted,
                })
            }
        }

        // Загружаем ключевые метрики
        metrics, err := wrapper.Db.Queries.GetClubKeyMetrics(ctx, club.ID)
        if err == nil {
            for _, metric := range metrics {
                sportClub.KeyMetrics = append(sportClub.KeyMetrics, sport_clubs.KeyMetric{
                    ID:          metric.ID,
                    ClubID:      metric.ClubID,
                    Name:        metric.Name,
                    Unit:        metric.Unit,
                    Description: metric.Description,
                    CreatedAt:   metric.CreatedAt,
                    UpdatedAt:   metric.UpdatedAt,
                    IsDeleted:   metric.IsDeleted,
                })
            }
        }

        response.Clubs = append(response.Clubs, sportClub)
    }

    ctx.JSON(http.StatusOK, response)
}

func GetSportClubHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
    id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Invalid club ID",
        })
        return
    }

    club, err := wrapper.Db.Queries.GetSportClubById(ctx, id)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{
            "message": "Club not found",
        })
        return
    }

    var sportClub sport_clubs.SportClub
    if err := mapper.Mapper{}.Map(&sportClub, club); err != nil {
        fmt.Printf("Error mapping club: %s\n", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to process club data",
        })
        return
    }

    // Загружаем расписание
    schedules, err := wrapper.Db.Queries.GetClubSchedules(ctx, club.ID)
    if err == nil {
        for _, schedule := range schedules {
            sportClub.Schedules = append(sportClub.Schedules, sport_clubs.Schedule{
                ID:        schedule.ID,
                ClubID:    schedule.ClubID,
                DayOfWeek: schedule.DayOfWeek,
                StartTime: schedule.StartTime,
                EndTime:   schedule.EndTime,
                Location:  schedule.Location,
                CreatedAt: schedule.CreatedAt,
                UpdatedAt: schedule.UpdatedAt,
                IsDeleted: schedule.IsDeleted,
            })
        }
    }

    // Загружаем ключевые метрики
    metrics, err := wrapper.Db.Queries.GetClubKeyMetrics(ctx, club.ID)
    if err == nil {
        for _, metric := range metrics {
            sportClub.KeyMetrics = append(sportClub.KeyMetrics, sport_clubs.KeyMetric{
                ID:          metric.ID,
                ClubID:      metric.ClubID,
                Name:        metric.Name,
                Unit:        metric.Unit,
                Description: metric.Description,
                CreatedAt:   metric.CreatedAt,
                UpdatedAt:   metric.UpdatedAt,
                IsDeleted:   metric.IsDeleted,
            })
        }
    }

    ctx.JSON(http.StatusOK, sportClub)
}

func UpdateSportClubHandler(ctx *gin.Context, wrapper *service_wrapper.Wrapper) {
    claimsRaw, exists := ctx.Get(middleware.ClaimsKey)
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "message": "Unauthorized",
        })
        return
    }

    userClaims := claimsRaw.(claims.UserClaims)
    id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Invalid club ID",
        })
        return
    }

    var request sport_clubs.UpdateSportClubRequest
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "Invalid request data",
        })
        return
    }

    // Проверяем существование клуба
    existingClub, err := wrapper.Db.Queries.GetSportClubById(ctx, id)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{
            "message": "Club not found",
        })
        return
    }

    // Проверяем права доступа (только создатель или админ)
    if userClaims.Role != "Admin" && existingClub.CreatedBy != userClaims.ID {
        ctx.JSON(http.StatusForbidden, gin.H{
            "message": "Access denied. Only club creator or admin can update the club",
        })
        return
    }

    // Обновляем основную информацию клуба
    updatedClub, err := wrapper.Db.Queries.UpdateSportClub(ctx, map[string]interface{}{
        "id":                id,
        "name":              request.Name,
        "description":       request.Description,
        "sport_type":        request.SportType,
        "total_slots":       request.TotalSlots,
        "reserved_slots":    request.ReservedSlots,
        "location":          request.Location,
        "min_hours_per_week": request.MinHoursPerWeek,
        "skill_level":       request.SkillLevel,
    })
    if err != nil {
        fmt.Printf("Error updating sport club: %s\n", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to update sport club",
        })
        return
    }

    // Управление расписанием
    err = updateClubSchedules(ctx, wrapper, id, request.Schedules)
    if err != nil {
        fmt.Printf("Error updating schedules: %s\n", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to update club schedules",
        })
        return
    }

    // Управление ключевыми метриками
    err = updateClubKeyMetrics(ctx, wrapper, id, request.KeyMetrics)
    if err != nil {
        fmt.Printf("Error updating key metrics: %s\n", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to update club key metrics",
        })
        return
    }

    // Получаем обновленные данные клуба
    finalClub, err := wrapper.Db.Queries.GetSportClubById(ctx, id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to fetch updated club data",
        })
        return
    }

    var response sport_clubs.SportClub
    if err := mapper.Mapper{}.Map(&response, finalClub); err != nil {
        fmt.Printf("Error mapping updated club: %s\n", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to process updated club data",
        })
        return
    }

    // Загружаем обновленное расписание и метрики
    schedules, _ := wrapper.Db.Queries.GetClubSchedules(ctx, id)
    for _, schedule := range schedules {
        response.Schedules = append(response.Schedules, sport_clubs.Schedule{
            ID:        schedule.ID,
            ClubID:    schedule.ClubID,
            DayOfWeek: schedule.DayOfWeek,
            StartTime: schedule.StartTime,
            EndTime:   schedule.EndTime,
            Location:  schedule.Location,
            CreatedAt: schedule.CreatedAt,
            UpdatedAt: schedule.UpdatedAt,
            IsDeleted: schedule.IsDeleted,
        })
    }

    metrics, _ := wrapper.Db.Queries.GetClubKeyMetrics(ctx, id)
    for _, metric := range metrics {
        response.KeyMetrics = append(response.KeyMetrics, sport_clubs.KeyMetric{
            ID:          metric.ID,
            ClubID:      metric.ClubID,
            Name:        metric.Name,
            Unit:        metric.Unit,
            Description: metric.Description,
            CreatedAt:   metric.CreatedAt,
            UpdatedAt:   metric.UpdatedAt,
            IsDeleted:   metric.IsDeleted,
        })
    }

    ctx.JSON(http.StatusOK, response)
}

// Вспомогательная функция для обновления расписания
func updateClubSchedules(ctx *gin.Context, wrapper *service_wrapper.Wrapper, clubID int64, schedules []sport_clubs.Schedule) error {
    // Получаем текущее расписание
    currentSchedules, err := wrapper.Db.Queries.GetClubSchedules(ctx, clubID)
    if err != nil {
        return err
    }

    // Создаем карту существующих расписаний для быстрого поиска
    existingSchedules := make(map[int64]bool)
    for _, schedule := range currentSchedules {
        existingSchedules[schedule.ID] = true
    }

    // Обновляем или создаем расписания
    for _, schedule := range schedules {
        if schedule.ID > 0 && existingSchedules[schedule.ID] {
            // Обновляем существующее расписание
            _, err := wrapper.Db.Queries.UpdateClubSchedule(ctx, map[string]interface{}{
                "id":         schedule.ID,
                "club_id":    clubID,
                "day_of_week": schedule.DayOfWeek,
                "start_time": schedule.StartTime,
                "end_time":   schedule.EndTime,
                "location":   schedule.Location,
            })
            if err != nil {
                return err
            }
        } else if schedule.ID == 0 {
            // Создаем новое расписание
            _, err := wrapper.Db.Queries.CreateClubSchedule(ctx, map[string]interface{}{
                "club_id":    clubID,
                "day_of_week": schedule.DayOfWeek,
                "start_time": schedule.StartTime,
                "end_time":   schedule.EndTime,
                "location":   schedule.Location,
            })
            if err != nil {
                return err
            }
        }
    }

    return nil
}

// Вспомогательная функция для обновления ключевых метрик
func updateClubKeyMetrics(ctx *gin.Context, wrapper *service_wrapper.Wrapper, clubID int64, metrics []sport_clubs.KeyMetric) error {
    // Получаем текущие метрики
    currentMetrics, err := wrapper.Db.Queries.GetClubKeyMetrics(ctx, clubID)
    if err != nil {
        return err
    }

    // Создаем карту существующих метрик для быстрого поиска
    existingMetrics := make(map[int64]bool)
    for _, metric := range currentMetrics {
        existingMetrics[metric.ID] = true
    }

    // Обновляем или создаем метрики
    for _, metric := range metrics {
        if metric.ID > 0 && existingMetrics[metric.ID] {
            // Обновляем существующую метрику
            _, err := wrapper.Db.Queries.UpdateClubKeyMetric(ctx, map[string]interface{}{
                "id":          metric.ID,
                "club_id":     clubID,
                "name":        metric.Name,
                "unit":        metric.Unit,
                "description": metric.Description,
            })
            if err != nil {
                return err
            }
        } else if metric.ID == 0 {
            // Создаем новую метрику
            _, err := wrapper.Db.Queries.CreateClubKeyMetric(ctx, map[string]interface{}{
                "club_id":     clubID,
                "name":        metric.Name,
                "unit":        metric.Unit,
                "description": metric.Description,
            })
            if err != nil {
                return err
            }
        }
    }

    return nil
}
