package version1

type Orderby struct {
	columns   string
	direction string
}

func Asc(columns string) Orderby {
	return Orderby{
		columns:   columns,
		direction: "ASC",
	}
}

func Desc(columns string) Orderby {
	return Orderby{
		columns:   columns,
		direction: "DESC",
	}
}
