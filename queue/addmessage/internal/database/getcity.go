package database

func (s StructDatabase) CheckCity(city string) (int, int, error) {

	var err error

	s.BaseOpen, err = s.DBJoin()
	if err != nil {
		return 0, 0, err
	}

	return s.sql(city)
}

func (s StructDatabase) sql(city string) (int, int, error) {

	var (
		reg int
		id  int
	)

	q := `SELECT region, id FROM city WHERE name_city LIKE ?`

	res, err := s.BaseOpen.Query(q, city)
	if err != nil {
		return 0, 0, err
	}
	defer res.Close()

	for res.Next() {
		if err = res.Scan(&reg, &id); err != nil {
			return 0, 0, err
		}
	}

	s.DBClose()
	return id, reg, err
}
