package charts

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// Repository chart database
type Repository struct {
	DB *sqlx.DB
}

// NewRepository create new repository
func NewRepository(DB *sqlx.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) getQuery(timeFormat string, timeOffset string, direction string) string {
	q := `
SELECT created_at 							     AS time,
       round(AVG(temp), 2)                       AS temp,
       c.name || ' [' || s.name || '] '          AS name,
       sensor_id
FROM sensors_data
         INNER JOIN sensors s on s.id = sensors_data.sensor_id
         INNER JOIN chips c on c.id = s.chip_id
WHERE created_at >= datetime('now', '%s')
GROUP BY strftime('%s', created_at), sensor_id
ORDER BY time %s
`

	return fmt.Sprintf(q, timeOffset, timeFormat, direction)
}

// BySeconds getting chart for last 1 minute
func (r *Repository) BySeconds() (*Chart, error) {
	var models []*Model

	q := r.getQuery("%Y-%m-%d %H:%M:%S", secondOffset, "ASC")

	if err := r.DB.Select(&models, q); err != nil {
		return nil, err
	}

	return mapToChart(models), nil
}

// ByMinutes getting chart for last 1 hour
func (r *Repository) ByMinutes() (*Chart, error) {
	var models []*Model

	q := r.getQuery("%Y-%m-%d %H:%M:00", minuteOffset, "ASC")

	if err := r.DB.Select(&models, q); err != nil {
		return nil, err
	}

	return mapToChart(models), nil
}

// ByHours getting chart for last 1 day
func (r *Repository) ByHours() (*Chart, error) {
	var models []*Model

	q := r.getQuery("%Y-%m-%d %H:00:00", hourOffset, "ASC")

	if err := r.DB.Select(&models, q); err != nil {
		return nil, err
	}

	return mapToChart(models), nil
}

// ByDays getting chart for last 1 month
func (r *Repository) ByDays() (*Chart, error) {
	var models []*Model

	q := r.getQuery("%Y-%m-%d 00:00:00", dayOffset, "ASC")

	if err := r.DB.Select(&models, q); err != nil {
		return nil, err
	}

	return mapToChart(models), nil
}

// ByMonth getting chart for last 1 year
func (r *Repository) ByMonth() (*Chart, error) {
	var models []*Model

	q := r.getQuery("%Y-%m-01 00:00:00", monthOffset, "ASC")

	if err := r.DB.Select(&models, q); err != nil {
		return nil, err
	}

	return mapToChart(models), nil
}

func (r *Repository) GetLast(typ string) (*Chart, error) {
	var models []*Model
	var offset string

	switch typ {
	case TypeSecond:
		offset = secondOffset
		break
	case TypeMinute:
		offset = minuteOffset
		break
	case TypeHour:
		offset = hourOffset
		break
	case TypeDay:
		offset = dayOffset
		break
	case TypeMonth:
		offset = monthOffset
		break
	default:
		offset = secondOffset
	}

	pQuery := r.getQuery("%Y-%m-%d %H:%M:%S", offset, "DESC")
	q := fmt.Sprintf("SELECT * FROM (%s) GROUP BY sensor_id", pQuery)

	if err := r.DB.Select(&models, q); err != nil {
		return nil, err
	}

	return mapToChart(models), nil
}
