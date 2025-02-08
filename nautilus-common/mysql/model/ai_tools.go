package model

type AITool struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement;comment:id"`
	Name        string `gorm:"type:varchar(255);not null;comment:tool name"`
	Description string `gorm:"type:text;not null;comment:tool description"`
	Strict      bool   `gorm:"type:char(1);default:0;comment:strict"`
	Parameters  string `gorm:"type:text;not null;comment:json parameters"`
	Role        string `gorm:"type:varchar(255);default:'*';comment:who can run this tools"`
}

func (*AITool) TableName() string {
	return "ai_tools"
}
