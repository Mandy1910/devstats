package main

import (
	"database/sql"
	lib "devstats"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type esRawCommit struct {
	Type                 string   `json:"type"`
	SHA                  string   `json:"sha"`
	EventID              int64    `json:"event_id"`
	AuthorName           string   `json:"author_name"`
	Message              string   `json:"message"`
	ActorLogin           string   `json:"actor_login"`
	RepoName             string   `json:"repo_name"`
	CreatedAt            string   `json:"time"`
	EncryptedEmail       string   `json:"encrypted_author_email"`
	AuthorEmail          string   `json:"author_email"`
	CommitterName        string   `json:"committer_name"`
	CommitterEmail       string   `json:"committer_email"`
	AuthorLogin          string   `json:"author_login"`
	CommitterLogin       string   `json:"committer_login"`
	Org                  string   `json:"org"`
	RepoGroup            string   `json:"repo_group"`
	RepoAlias            string   `json:"repo_alias"`
	ActorName            string   `json:"actor_name"`
	ActorCountryCode     string   `json:"actor_country_code"`
	ActorGender          string   `json:"actor_gender"`
	ActorGenderProb      *float64 `json:"actor_gender_prob"`
	ActorTZ              string   `json:"actor_tz"`
	ActorTZOffset        *int     `json:"actor_tz_offset"`
	ActorCountry         string   `json:"actor_country"`
	AuthorCountryCode    string   `json:"author_country_code"`
	AuthorGender         string   `json:"author_gender"`
	AuthorGenderProb     *float64 `json:"author_gender_prob"`
	AuthorTZ             string   `json:"author_tz"`
	AuthorTZOffset       *int     `json:"author_tz_offset"`
	AuthorCountry        string   `json:"author_country"`
	CommitterCountryCode string   `json:"committer_country_code"`
	CommitterGender      string   `json:"committer_gender"`
	CommitterGenderProb  *float64 `json:"committer_gender_prob"`
	CommitterTZ          string   `json:"committer_tz"`
	CommitterTZOffset    *int     `json:"committer_tz_offset"`
	CommitterCountry     string   `json:"committer_country"`
	ActorCompany         string   `json:"actor_company"`
	AuthorCompany        string   `json:"author_company"`
	Committer            string   `json:"committer_company"`
	Size                 int      `json:"size"`
}

type esRawIssue struct {
	Type                string   `json:"type"`
	ID                  int64    `json:"id"`
	EventID             int64    `json:"event_id"`
	EventCreatedAt      string   `json:"time"`
	CreatedAt           string   `json:"created_at"`
	Body                string   `json:"body"`
	ClosedAt            *string  `json:"closed_at"`
	Comments            int      `json:"comments"`
	Locked              bool     `json:"locked"`
	Number              int      `json:"number"`
	State               string   `json:"state"`
	Title               string   `json:"title"`
	UpdatedAt           string   `json:"updated_at"`
	IsPR                bool     `json:"is_pr"`
	EventType           string   `json:"event_type"`
	MilestoneNumber     *int     `json:"milestone_number"`
	MilestoneState      string   `json:"milestone_state"`
	MilestoneTitle      string   `json:"milestone_title"`
	AssigneeLogin       string   `json:"assignee_login"`
	AssigneeName        string   `json:"assignee_name"`
	AssigneeCountryCode string   `json:"assignee_country_code"`
	AssigneeGender      string   `json:"assignee_gender"`
	AssigneeGenderProb  *float64 `json:"assignee_gender_prob"`
	AssigneeTZ          string   `json:"assignee_tz"`
	AssigneeTZOffset    *int     `json:"assignee_tz_offset"`
	AssigneeCountry     string   `json:"assignee_country"`
	ActorLogin          string   `json:"actor_login"`
	ActorName           string   `json:"actor_name"`
	ActorCountryCode    string   `json:"actor_country_code"`
	ActorGender         string   `json:"actor_gender"`
	ActorGenderProb     *float64 `json:"actor_gender_prob"`
	ActorTZ             string   `json:"actor_tz"`
	ActorTZOffset       *int     `json:"actor_tz_offset"`
	ActorCountry        string   `json:"actor_country"`
	UserLogin           string   `json:"user_login"`
	UserName            string   `json:"user_name"`
	UserCountryCode     string   `json:"user_country_code"`
	UserGender          string   `json:"user_gender"`
	UserGenderProb      *float64 `json:"user_gender_prob"`
	UserTZ              string   `json:"user_tz"`
	UserTZOffset        *int     `json:"user_tz_offset"`
	UserCountry         string   `json:"user_country"`
}

func generateRawES(ch chan struct{}, ctx *lib.Ctx, con *sql.DB, es *lib.ES, dtf, dtt time.Time, sqls map[string]string) {
	if ctx.Debug > 0 {
		lib.Printf("Working on %v - %v\n", dtf, dtt)
	}

	// Replace dates
	sFrom := lib.ToYMDHMSDate(dtf)
	sTo := lib.ToYMDHMSDate(dtt)

	// ES bulk inserts
	bulkDel, bulkAdd := es.Bulks()

	// Commits
	sql := strings.Replace(sqls["commits"], "{{from}}", sFrom, -1)
	sql = strings.Replace(sql, "{{to}}", sTo, -1)

	// Execute query
	rows := lib.QuerySQLWithErr(con, ctx, sql)
	defer func() { lib.FatalOnError(rows.Close()) }()

	var (
		commit    esRawCommit
		createdAt time.Time
	)
	shas := make(map[string]struct{})
	commit.Type = "commit"
	nCommits := 0
	for rows.Next() {
		lib.FatalOnError(
			rows.Scan(
				&commit.SHA,
				&commit.EventID,
				&commit.AuthorName,
				&commit.Message,
				&commit.ActorLogin,
				&commit.RepoName,
				&createdAt,
				&commit.EncryptedEmail,
				&commit.AuthorEmail,
				&commit.CommitterName,
				&commit.CommitterEmail,
				&commit.AuthorLogin,
				&commit.CommitterLogin,
				&commit.Org,
				&commit.RepoGroup,
				&commit.RepoAlias,
				&commit.ActorName,
				&commit.ActorCountryCode,
				&commit.ActorGender,
				&commit.ActorGenderProb,
				&commit.ActorTZ,
				&commit.ActorTZOffset,
				&commit.ActorCountry,
				&commit.AuthorCountryCode,
				&commit.AuthorGender,
				&commit.AuthorGenderProb,
				&commit.AuthorTZ,
				&commit.AuthorTZOffset,
				&commit.AuthorCountry,
				&commit.CommitterCountryCode,
				&commit.CommitterGender,
				&commit.CommitterGenderProb,
				&commit.CommitterTZ,
				&commit.CommitterTZOffset,
				&commit.CommitterCountry,
				&commit.ActorCompany,
				&commit.AuthorCompany,
				&commit.Committer,
				&commit.Size,
			),
		)
		nCommits++
		commit.CreatedAt = lib.ToESDate(createdAt)
		commit.Message = lib.TruncToBytes(commit.Message, 0x400)
		shas[commit.SHA] = struct{}{}
		es.AddBulksItemsI(ctx, bulkDel, bulkAdd, commit, lib.HashArray([]interface{}{commit.Type, commit.SHA, commit.EventID}))
		if nCommits%10000 == 0 {
			// Bulk insert to ES
			es.ExecuteBulks(ctx, bulkDel, bulkAdd)
		}
	}
	lib.FatalOnError(rows.Err())

	// Issues
	sql = strings.Replace(sqls["issues"], "{{from}}", sFrom, -1)
	sql = strings.Replace(sql, "{{to}}", sTo, -1)

	// Execute query
	rows = lib.QuerySQLWithErr(con, ctx, sql)
	defer func() { lib.FatalOnError(rows.Close()) }()

	var (
		issue          esRawIssue
		eventCreatedAt time.Time
		closedAt       *time.Time
		updatedAt      time.Time
	)
	ids := make(map[int64]struct{})
	issue.Type = "issue"
	nIssues := 0
	for rows.Next() {
		lib.FatalOnError(
			rows.Scan(
				&issue.ID,
				&issue.EventID,
				&eventCreatedAt,
				&createdAt,
				&issue.Body,
				&closedAt,
				&issue.Comments,
				&issue.Locked,
				&issue.Number,
				&issue.State,
				&issue.Title,
				&updatedAt,
				&issue.IsPR,
				&issue.EventType,
				&issue.MilestoneNumber,
				&issue.MilestoneState,
				&issue.MilestoneTitle,
				&issue.AssigneeLogin,
				&issue.AssigneeName,
				&issue.AssigneeCountryCode,
				&issue.AssigneeGender,
				&issue.AssigneeGenderProb,
				&issue.AssigneeTZ,
				&issue.AssigneeTZOffset,
				&issue.AssigneeCountry,
				&issue.ActorLogin,
				&issue.ActorName,
				&issue.ActorCountryCode,
				&issue.ActorGender,
				&issue.ActorGenderProb,
				&issue.ActorTZ,
				&issue.ActorTZOffset,
				&issue.ActorCountry,
				&issue.UserLogin,
				&issue.UserName,
				&issue.UserCountryCode,
				&issue.UserGender,
				&issue.UserGenderProb,
				&issue.UserTZ,
				&issue.UserTZOffset,
				&issue.UserCountry,
			),
		)
		nIssues++
		issue.CreatedAt = lib.ToESDate(createdAt)
		issue.EventCreatedAt = lib.ToESDate(eventCreatedAt)
		issue.UpdatedAt = lib.ToESDate(updatedAt)
		issue.Body = lib.TruncToBytes(issue.Body, 0x400)
		if closedAt != nil {
			tm := lib.ToESDate(*closedAt)
			issue.ClosedAt = &tm
		} else {
			issue.ClosedAt = nil
		}
		ids[issue.ID] = struct{}{}
		es.AddBulksItemsI(ctx, bulkDel, bulkAdd, issue, lib.HashArray([]interface{}{issue.Type, issue.ID, issue.EventID}))
		if nIssues%10000 == 0 {
			// Bulk insert to ES
			es.ExecuteBulks(ctx, bulkDel, bulkAdd)
		}
	}
	lib.FatalOnError(rows.Err())

	// Bulk insert to ES
	es.ExecuteBulks(ctx, bulkDel, bulkAdd)

	if ctx.Debug > 0 {
		lib.Printf(
			"%v - %v: %d commits (%d unique SHAs), %d issue events (%d unique issues)\n",
			sFrom, sTo, nCommits, len(shas), nIssues, len(ids),
		)
	}

	if ch != nil {
		ch <- struct{}{}
	}
}

// gha2es - main working function
func gha2es(args []string) {
	var (
		ctx      lib.Ctx
		err      error
		hourFrom int
		hourTo   int
		dFrom    time.Time
		dTo      time.Time
	)

	// Environment context parse
	ctx.Init()
	if !ctx.UseES {
		return
	}
	// Connect to ElasticSearch
	es := lib.ESConn(&ctx, "d_raw_")
	// Create index
	exists := es.IndexExists(&ctx)
	if !exists {
		es.CreateIndex(&ctx, true)
	}

	// Connect to Postgres DB
	con := lib.PgConn(&ctx)
	defer func() { lib.FatalOnError(con.Close()) }()

	// Get raw commits to ES SQL
	sqls := make(map[string]string)
	dataPrefix := lib.DataDir
	if ctx.Local {
		dataPrefix = "./"
	}
	data := [][2]string{
		{"commits", "util_sql/es_raw_commits.sql"},
		{"issues", "util_sql/es_raw_issues.sql"},
	}
	for _, row := range data {
		bytes, err := lib.ReadFile(
			&ctx,
			dataPrefix+row[1],
		)
		lib.FatalOnError(err)
		sqls[row[0]] = string(bytes)
	}

	// Current date
	now := time.Now()
	startD, startH, endD, endH := args[0], args[1], args[2], args[3]

	// Parse from day & hour
	if strings.ToLower(startH) == lib.Now {
		hourFrom = now.Hour()
	} else {
		hourFrom, err = strconv.Atoi(startH)
		lib.FatalOnError(err)
	}

	if strings.ToLower(startD) == lib.Today {
		dFrom = lib.DayStart(now).Add(time.Duration(hourFrom) * time.Hour)
	} else {
		dFrom, err = time.Parse(
			time.RFC3339,
			fmt.Sprintf("%sT%02d:00:00+00:00", startD, hourFrom),
		)
		lib.FatalOnError(err)
	}

	// Parse to day & hour
	if strings.ToLower(endH) == lib.Now {
		hourTo = now.Hour()
	} else {
		hourTo, err = strconv.Atoi(endH)
		lib.FatalOnError(err)
	}

	if strings.ToLower(endD) == lib.Today {
		dTo = lib.DayStart(now).Add(time.Duration(hourTo) * time.Hour)
	} else {
		dTo, err = time.Parse(
			time.RFC3339,
			fmt.Sprintf("%sT%02d:00:00+00:00", endD, hourTo),
		)
		lib.FatalOnError(err)
	}

	// Get number of CPUs available and optimal time window for threads
	thrN := lib.GetThreadsNum(&ctx)
	hours := int(dTo.Sub(dFrom).Hours()) / thrN
	if hours < 1 {
		hours = 1
	}
	lib.Printf("gha2es.go: Running (%v CPUs): %v - %v, interval %dh\n", thrN, dFrom, dTo, hours)

	dt := dFrom
	dtN := dt
	if thrN > 1 {
		ch := make(chan struct{})
		nThreads := 0
		for dt.Before(dTo) || dt.Equal(dTo) {
			dtN = dt.Add(time.Hour * time.Duration(hours))
			go generateRawES(ch, &ctx, con, es, dt, dtN, sqls)
			dt = dtN
			nThreads++
			if nThreads == thrN {
				<-ch
				nThreads--
			}
		}
		lib.Printf("Final threads join\n")
		for nThreads > 0 {
			<-ch
			nThreads--
		}
	} else {
		lib.Printf("Using single threaded version\n")
		for dt.Before(dTo) || dt.Equal(dTo) {
			dtN = dt.Add(time.Hour * time.Duration(hours))
			generateRawES(nil, &ctx, con, es, dt, dtN, sqls)
			dt = dtN
		}
	}
	// Finished
	lib.Printf("All done.\n")
}

func main() {
	dtStart := time.Now()
	// Required args
	if len(os.Args) < 4 {
		lib.Printf("Arguments required: date_from_YYYY-MM-DD hour_from_HH date_to_YYYY-MM-DD hour_to_HH\n")
		os.Exit(1)
	}
	gha2es(os.Args[1:])
	dtEnd := time.Now()
	lib.Printf("Time: %v\n", dtEnd.Sub(dtStart))
}
