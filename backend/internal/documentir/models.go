package documentir

type BlockType string

const (
	TypeHeading   BlockType = "heading"
	TypeParagraph BlockType = "paragraph"
	TypeEquation  BlockType = "equation"
	TypeList      BlockType = "list"
	TypeTable     BlockType = "table"
	TypeImage     BlockType = "image"
	TypeDiagram   BlockType = "diagram"
	TypeGraph     BlockType = "graph"
	TypeCode      BlockType = "code"
	TypeQuote     BlockType = "quote"
	TypePageBreak BlockType = "page_break"
)

type Block struct {
	ID         string                 `json:"id"`
	Type       BlockType              `json:"type"`
	Confidence float64                `json:"confidence"`
	Content    map[string]interface{} `json:"content"` // Typed based on BlockType in actual usage
}

type Page struct {
	PageNumber int     `json:"page_number"`
	Blocks     []Block `json:"blocks"`
}

type Document struct {
	Pages []Page `json:"pages"`
}
