package model

type Obfuscation struct {
	ID        string `db:"id"`
	OrgID     string `db:"org_id"`
	RuleName  string `db:"rule_name"`
	RuleType  string `db:"rule_type"`
	RuleDef   []byte `db:"rule_def"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
	Deleted   bool   `db:"deleted"`
	Disabled  bool   `db:"disabled"`
}

func (obfuscation Obfuscation) GetAllColumns() []any {
	return []any{obfuscation.OrgID, obfuscation.RuleName, obfuscation.RuleType, obfuscation.RuleDef, obfuscation.CreatedAt, obfuscation.UpdatedAt, obfuscation.Deleted, obfuscation.Disabled}
}
