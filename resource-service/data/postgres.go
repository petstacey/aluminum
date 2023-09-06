package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

const (
	dbTimeout = 5 * time.Second
)

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		DB: db,
	}
}

func (r *PostgresRepository) CreateResource(resource *Resource) error {
	query := `INSERT INTO resources
			  (id, name, email, type_id, job_title_id, workgroup_id, location_id, manager_id)
			  VALUES ($1, $2, $3,
				(SELECT typ.id FROM employment_types typ WHERE type = $4),
			    (SELECT t.id FROM job_titles t WHERE title = $5),
			    (SELECT w.id FROM workgroups w WHERE w.name = $6),
			    (SELECT l.id FROM locations l WHERE l.name = $7),
			    (SELECT m.id FROM resources m WHERE m.name = $8)
			  ) RETURNING active`
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	args := []interface{}{
		resource.ID,
		resource.Name,
		resource.Email,
		resource.Type,
		resource.JobTitle,
		resource.Workgroup,
		resource.Location,
		resource.Manager,
	}
	return r.DB.QueryRowContext(ctx, query, args...).Scan(&resource.Active)
}

func (r *PostgresRepository) GetResource(id int64) (*Resource, error) {
	query := `SELECT r.id, r.name, r.email, typ.type, j.title, w.name, l.name, m.name, r.active
	          FROM resources r
			  	JOIN employment_types typ ON (typ.id = r.type_id)
				JOIN job_titles j ON (j.id = r.job_title_id)
				JOIN workgroups w ON ( w.id = r.workgroup_id)
				JOIN locations l ON (l.id = r.location_id)
				JOIN resources m ON ( m.id = r.manager_id)
			  WHERE r.id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var res Resource
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&res.ID,
		&res.Name,
		&res.Email,
		&res.Type,
		&res.JobTitle,
		&res.Workgroup,
		&res.Location,
		&res.Manager,
		&res.Active,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *PostgresRepository) GetResources(name string, titles, types, workgroups, locations, managers []string, filters Filters) ([]*Resource, Metadata, error) {
	query := fmt.Sprintf(`
        SELECT count(*) OVER(), r.id, r.name, r.email, typ.type, t.title, w.name, l.name, m.name AS manager, r.active
        FROM (((((resources r
            JOIN job_titles t ON r.job_title_id = t.id)
			JOIN employment_types typ ON r.type_id = typ.id)
            JOIN workgroups w ON w.id = r.workgroup_id)
			JOIN locations l ON l.id = r.location_id)
            JOIN resources m ON r.manager_id = m.id)
        WHERE (r.name = $1 OR $1 = '')
		AND (t.title = ANY($2) OR $2 = '{}')
		AND (typ.type = ANY($3) OR $3 = '{}')
		AND (w.name = ANY($4) OR $4 = '{}')
		AND (l.name = ANY($5) OR $5 = '{}')
        AND (m.name = ANY($6) OR $6 = '{}')
        ORDER BY %s %s, r.id ASC
        LIMIT $7 OFFSET $8`, fmt.Sprintf("r.%s", filters.sortColumn()), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	args := []interface{}{
		name,
		pq.Array(titles),
		pq.Array(types),
		pq.Array(workgroups),
		pq.Array(locations),
		pq.Array(managers),
		filters.limit(),
		filters.offset(),
	}
	totalRecords := 0
	resources := []*Resource{}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var res Resource
		err := rows.Scan(
			&totalRecords,
			&res.ID,
			&res.Name,
			&res.Email,
			&res.Type,
			&res.JobTitle,
			&res.Workgroup,
			&res.Location,
			&res.Manager,
			&res.Active,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		resources = append(resources, &res)
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return resources, metadata, nil
}

func (r *PostgresRepository) UpdateResource(resource *Resource) error {
	query := `UPDATE resources 
	          SET name = $1, email = $2,
	          type_id = (SELECT typ.id FROM employment_types typ WHERE type = $3),
			  job_title_id = (SELECT j.id FROM job_titles j WHERE j.title = $4),
			  workgroup_id = (SELECT w.id FROM workgroups w WHERE w.name = $5),
			  location_id = (SELECT l.id FROM locations l WHERE l.name = $6),
			  manager_id = (SELECT m.id FROM resources m WHERE m.name = $7),
			  active = $8
			  WHERE id = $9`
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	args := []interface{}{
		resource.Name,
		resource.Email,
		resource.Type,
		resource.JobTitle,
		resource.Workgroup,
		resource.Location,
		resource.Manager,
		resource.Active,
		resource.ID,
	}
	err := r.DB.QueryRowContext(ctx, query, args...).Scan()
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil
		default:
			return err
		}
	}
	return nil
}
