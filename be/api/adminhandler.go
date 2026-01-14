package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rprj/be/dblayer"
	"strings"
)

type DashboardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	// Define fields for dashboard data as needed
	UsersCount  int                       `json:"users_count"`
	UsersStats  map[string]int            `json:"users_stats"`
	GroupsCount int                       `json:"groups_count"`
	ObjectStats map[string]map[string]int `json:"object_stats"` // e.g., {"Project": {"count": 10, "deleted_count": 2}, ...}
}

// DashboardHandler godoc
// @Summary Admin Dashboard
// @Description Returns admin dashboard data
// @Tags admin
// @Produce json
// @Success 200 {object} DashboardResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /admin/dashboard [get]
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := GetClaimsFromRequest(r)
	if err != nil {
		RespondSimpleError(w, ErrUnauthorized, "Unauthorized", http.StatusUnauthorized)
		return
	}

	dbContext := &dblayer.DBContext{
		UserID:   claims["user_id"],
		GroupIDs: strings.Split(claims["groups"], ","),
		Schema:   dblayer.DbSchema,
	}

	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = false

	var response DashboardResponse
	response.Success = true
	response.Message = "Dashboard data retrieved successfully"

	// Get user statistics
	err = userStatistics(repo, &response)
	if err != nil {
		response.Success = false
		response.Message = "Failed to get dashboard data"
		log.Printf("DashboardHandler: error getting user statistics: %v\n", err)
	}
	// Get group statistics
	err = groupStatistics(repo, &response)
	if err != nil {
		response.Success = false
		response.Message = "Failed to get dashboard data"
		log.Printf("DashboardHandler: error getting group statistics: %v\n", err)
	}

	// For each object type, gather relevant stats
	response.ObjectStats = make(map[string]map[string]int)
	for _, className := range dblayer.Factory.GetAllClassNames() {

		objectStats, err := objectStatistics(className, repo, &response)
		if err != nil {
			response.Success = false
			response.Message = "Failed to get dashboard data"
			log.Printf("DashboardHandler: error getting statistics for %s: %v\n", className, err)
			break
		}

		response.ObjectStats[className] = objectStats
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func userStatistics(repo *dblayer.DBRepository, response *DashboardResponse) error {
	// Count users
	results := repo.Select("DBObject", "select count(*) as num from "+repo.DbContext.Schema+"_"+"users")
	fmt.Println("userStatistics: users count query result:", len(results))
	if len(results) == 0 {
		response.Success = false
		response.Message = "Failed to get dashboard data"
		log.Print("userStatistics: error getting users count: no results")
		return fmt.Errorf("no results for users count")
	}
	fmt.Println("userStatistics: users count =", results[0].GetValue("num"))
	tmpStr := results[0].GetValue("num").(string)
	fmt.Print("userStatistics: users count string =", tmpStr, "\n")
	_, err := fmt.Sscanf(tmpStr, "%d", &response.UsersCount)
	if err != nil {
		response.Success = false
		response.Message = "Failed to parse users count"
		log.Print("userStatistics: error parsing users count:", err)
		return err
	}

	// **** users stats
	userStats := make(map[string]int)

	// Users active last 24h
	queryActiveLastDay := "select count(distinct user_id) as num from " + repo.DbContext.Schema + "_" + "oauth_tokens where user_id in (select id from rra_users) and created_at >= NOW() - INTERVAL 1 day"
	fmt.Println("userStatistics: active users last 24h query:", queryActiveLastDay)
	results = repo.Select("DBObject", queryActiveLastDay)
	if len(results) == 1 {
		fmt.Println("userStatistics: active users last 24h =", results[0].GetValue("num"))
		tmpStr := results[0].GetValue("num").(string)
		tmpInt := 0
		_, err = fmt.Sscanf(tmpStr, "%d", &tmpInt)
		if err != nil {
			response.Success = false
			response.Message = "Failed to parse active users last 24h count"
			log.Print("userStatistics: error parsing active users last 24h count:", err)
			return err
		}
		userStats["active_last_24h"] = tmpInt
	}
	// Users active last 7 days
	queryActiveLastWeek := "select count(distinct user_id) as num from " + repo.DbContext.Schema + "_" + "oauth_tokens where user_id in (select id from rra_users) and created_at >= NOW() - INTERVAL 7 day"
	fmt.Println("userStatistics: active users last 7 days query:", queryActiveLastWeek)
	results = repo.Select("DBObject", queryActiveLastWeek)
	if len(results) == 1 {
		fmt.Println("userStatistics: active users last 7 days =", results[0].GetValue("num"))
		tmpStr := results[0].GetValue("num").(string)
		tmpInt := 0
		_, err = fmt.Sscanf(tmpStr, "%d", &tmpInt)
		if err != nil {
			response.Success = false
			response.Message = "Failed to parse active users last 7 days count"
			log.Print("userStatistics: error parsing active users last 7 days count:", err)
			return err
		}
		userStats["active_last_7_days"] = tmpInt
	}
	// Users active last 30 days
	queryActiveLastMonth := "select count(distinct user_id) as num from " + repo.DbContext.Schema + "_" + "oauth_tokens where user_id in (select id from rra_users) and created_at >= NOW() - INTERVAL 30 day"
	fmt.Println("userStatistics: active users last 30 days query:", queryActiveLastMonth)
	results = repo.Select("DBObject", queryActiveLastMonth)
	if len(results) == 1 {
		fmt.Println("userStatistics: active users last 30 days =", results[0].GetValue("num"))
		tmpStr := results[0].GetValue("num").(string)
		tmpInt := 0
		_, err = fmt.Sscanf(tmpStr, "%d", &tmpInt)
		if err != nil {
			response.Success = false
			response.Message = "Failed to parse active users last 30 days count"
			log.Print("userStatistics: error parsing active users last 30 days count:", err)
			return err
		}
		userStats["active_last_30_days"] = tmpInt
	}
	response.UsersStats = userStats

	// select user_id, max(created_at) from rra_oauth_tokens where user_id in (select id from rra_users) group by user_id;

	return nil
}

func groupStatistics(repo *dblayer.DBRepository, response *DashboardResponse) error {
	// **** Count groups
	results := repo.Select("DBObject", "select count(*) as num from "+repo.DbContext.Schema+"_"+"groups")
	fmt.Println("groupStatistics: groups count query result:", len(results))
	if len(results) == 0 {
		response.Success = false
		response.Message = "Failed to get dashboard data"
		log.Print("groupStatistics: error getting groups count: no results")
		return fmt.Errorf("no results for groups count")
	}
	fmt.Println("groupStatistics: groups count =", results[0].GetValue("num"))
	tmpStr := results[0].GetValue("num").(string)
	fmt.Print("groupStatistics: groups count string =", tmpStr, "\n")
	// Convert string to int
	_, err := fmt.Sscanf(tmpStr, "%d", &response.GroupsCount)
	if err != nil {
		response.Success = false
		response.Message = "Failed to parse groups count"
		log.Print("groupStatistics: error parsing groups count:", err)
		return err
	}

	return nil
}

func objectStatistics(className string, repo *dblayer.DBRepository, response *DashboardResponse) (map[string]int, error) {
	dbe := dblayer.Factory.GetInstanceByClassName(className)
	if !dbe.IsDBObject() {
		return nil, nil
	}

	var tmpInt int
	var tmpStr string
	var err error

	objectStats := make(map[string]int)
	objectStats["count"] = 0
	tableName := dbe.GetTableName()

	// Count objects of this type
	results := repo.Select("DBObject", "select count(*) as num from "+repo.DbContext.Schema+"_"+tableName)
	if len(results) == 1 {
		fmt.Printf("objectStatistics: %s count = %v\n", className, results[0].GetValue("num"))
		tmpStr = results[0].GetValue("num").(string)
		_, err = fmt.Sscanf(tmpStr, "%d", &tmpInt)
		if err != nil {
			log.Printf("objectStatistics: error parsing %s count: %v\n", className, err)
			return nil, err
		}
		objectStats["count"] = tmpInt
	}
	// Count deleted objects of this type
	results = repo.Select("DBObject", "select count(*) as num from "+repo.DbContext.Schema+"_"+tableName+" where deleted_date is not null")
	if len(results) == 1 {
		fmt.Printf("objectStatistics: %s deleted count = %v\n", className, results[0].GetValue("num"))
		tmpStr = results[0].GetValue("num").(string)
		_, err = fmt.Sscanf(tmpStr, "%d", &tmpInt)
		if err != nil {
			return nil, err
		}
		objectStats["deleted_count"] = tmpInt
	}
	// Count created last week and modified last week
	results = repo.Select("DBObject", "select count(*) as num from "+repo.DbContext.Schema+"_"+tableName+" where creation_date >= NOW() - INTERVAL 7 day")
	if len(results) == 1 {
		fmt.Printf("objectStatistics: %s created last week count = %v\n", className, results[0].GetValue("num"))
		tmpStr = results[0].GetValue("num").(string)
		_, err = fmt.Sscanf(tmpStr, "%d", &tmpInt)
		if err != nil {
			return nil, err
		}
		objectStats["created_last_week"] = tmpInt
	}
	results = repo.Select("DBObject", "select count(*) as num from "+repo.DbContext.Schema+"_"+tableName+" where last_modify_date >= NOW() - INTERVAL 7 day")
	if len(results) == 1 {
		fmt.Printf("objectStatistics: %s modified last week count = %v\n", className, results[0].GetValue("num"))
		tmpStr = results[0].GetValue("num").(string)
		_, err = fmt.Sscanf(tmpStr, "%d", &tmpInt)
		if err != nil {
			return nil, err
		}
		objectStats["modified_last_week"] = tmpInt
	}
	return objectStats, nil
}
