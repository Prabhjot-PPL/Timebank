package dto

type UserDetails struct {
	Id            int      `json:"id"`
	Username      string   `json:"username"`
	Email         string   `json:"email"`
	Password      string   `json:"password"`
	SkillsOffered []string `json:"skills_offered"`
	SkillsNeeded  []string `json:"skills_needed"`
	BalanceHours  float64  `json:"balance_hours"`
}

type HelperDetails struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	SkillOffered string `json:"skill_offered"`
}
