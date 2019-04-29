package model

import "fmt"

// Region save region
type Region struct {
	ID        int
	NumRegion int
	Name      string
	IsTown    bool
}

// SelectRegions get regions
func (d *Database) SelectRegions() ([]Region, error) {
	if d.err != nil {
		return nil, d.err
	}
	rows, err := d.db.Query(sqlGetRegions)
	if err != nil {
		return nil, fmt.Errorf("[db] query %v", err)
	}
	regions := []Region{}

	for rows.Next() {
		r := Region{}

		if err := rows.Scan(&r.ID, &r.NumRegion, &r.Name, &r.IsTown); err != nil {
			return nil, fmt.Errorf("[db] scan %v", err)
		}
		regions = append(regions, r)
	}
	return regions, nil
}
