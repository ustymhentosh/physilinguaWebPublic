package problem

type Problem struct {
	Status string
	Id     string
}

func GetFilteredIdList(problems_list []Problem, filter string) []string {
	var filtered_problems []string
	for _, i := range problems_list {
		if i.Status == filter {
			filtered_problems = append(filtered_problems, i.Id)
		}
	}
	return filtered_problems
}
